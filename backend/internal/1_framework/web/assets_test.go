package web

import (
	"crypto/sha256"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestWebTemplatesParse(t *testing.T) {
	t.Parallel()
	if _, err := template.ParseGlob("../../../web/*.html"); err != nil {
		t.Fatalf("parse web templates: %v", err)
	}
}

func TestApplicationScriptsAvoidJQueryEventShortcuts(t *testing.T) {
	t.Parallel()
	scripts, err := filepath.Glob("../../../web/js/*.js")
	if err != nil {
		t.Fatalf("list application scripts: %v", err)
	}
	eventShortcut := regexp.MustCompile(`\.(click|change|submit|hover|focus|blur|keydown|keyup|keypress|mousedown|mouseup|mouseenter|mouseleave|mousemove|resize|scroll)\s*\(`)
	for _, script := range scripts {
		content, readErr := os.ReadFile(script) //nolint:gosec
		if readErr != nil {
			t.Fatalf("read script %s: %v", script, readErr)
		}
		if shortcut := eventShortcut.Find(content); shortcut != nil {
			t.Errorf("%s contains unsupported Cash event shortcut %q", filepath.Base(script), shortcut)
		}
	}
}

func TestApplicationScriptsPassDOMElementsToBootstrapComponents(t *testing.T) {
	t.Parallel()
	scripts, err := filepath.Glob("../../../web/js/*.js")
	if err != nil {
		t.Fatalf("list application scripts: %v", err)
	}
	cashBootstrapComponent := regexp.MustCompile(`new\s+bootstrap\.[A-Za-z]+\s*\(\s*\$\(`)
	for _, script := range scripts {
		content, readErr := os.ReadFile(script) //nolint:gosec
		if readErr != nil {
			t.Fatalf("read script %s: %v", script, readErr)
		}
		if match := cashBootstrapComponent.Find(content); match != nil {
			t.Errorf("%s passes a Cash collection to a Bootstrap component: %q", filepath.Base(script), match)
		}
	}
}

func TestAdminProductReadsInputsWithinProductElement(t *testing.T) {
	t.Parallel()
	content, err := os.ReadFile("../../../web/js/admin_product.js")
	if err != nil {
		t.Fatalf("read admin product script: %v", err)
	}
	script := string(content)
	for _, expected := range []string{
		"document.getElementsByName(String(janCode))[0]",
		`product.querySelector('input[type="checkbox"]')`,
		`product.querySelector('input[name="stock"]')`,
	} {
		if !strings.Contains(script, expected) {
			t.Errorf("admin product script does not contain %q", expected)
		}
	}
	if strings.Contains(script, "div[name=${janCode}]") {
		t.Error("admin product script contains an unquoted dynamic attribute selector")
	}
}

func TestRepositoryAggregatesPreparingOrdersWithMapAPI(t *testing.T) {
	t.Parallel()
	content, err := os.ReadFile("../../../web/js/repository.js")
	if err != nil {
		t.Fatalf("read repository script: %v", err)
	}
	script := string(content)
	for _, expected := range []string{
		"order.order_status === orderStatusPreparing",
		"placeJANCountMap.get(order.jan_code)",
		"placeJANCountMap.set(order.jan_code, currentCount + order.count)",
		"placeJANCountMap.forEach(function (count, janCode)",
	} {
		if !strings.Contains(script, expected) {
			t.Errorf("repository script does not contain %q", expected)
		}
	}
	if strings.Contains(script, "$.each(placeJANCountMap") {
		t.Error("repository script iterates a Map with the Cash object iterator")
	}
}

func TestWebAssetsUseSameOriginPaths(t *testing.T) {
	t.Parallel()
	pages, err := filepath.Glob("../../../web/*.html")
	if err != nil {
		t.Fatalf("list templates: %v", err)
	}
	absoluteAssetPrefixes := []string{
		`{{ .URL }}/css/`,
		`{{ .URL }}/fontawesome/`,
		`{{ .URL }}/image/`,
		`{{ .URL }}/js/`,
		`{{ .URL }}/sound/`,
	}
	for _, page := range pages {
		content, readErr := os.ReadFile(page) //nolint:gosec
		if readErr != nil {
			t.Fatalf("read template %s: %v", page, readErr)
		}
		for _, prefix := range absoluteAssetPrefixes {
			if strings.Contains(string(content), prefix) {
				t.Errorf("%s contains cross-origin asset prefix %q", filepath.Base(page), prefix)
			}
		}
	}
}

func TestCDNAssetsArePinnedAndIntegrityProtected(t *testing.T) {
	t.Parallel()
	expectedIntegrityByURL := map[string]string{
		"https://cdn.jsdelivr.net/npm/alpinejs@3.17.4/dist/cdn.min.js":                 "sha384-5/joNqFnRyVWzXp99bHot6RHG+EksGp+USSgZwPar7T9SD9PKKER37n/8bXBAZGd",
		"https://cdn.jsdelivr.net/npm/bootstrap@5.3.1/dist/css/bootstrap.min.css":      "sha384-4bw+/aepP/YC94hEpVNVgiZdgIC5+VKNBQNGCHeKRQN+PtmoHDEXuppvnDJzQIu9",
		"https://cdn.jsdelivr.net/npm/bootstrap@5.3.1/dist/js/bootstrap.bundle.min.js": "sha384-HwwvtgBNo3bZJJLYd8oVXjrBZt8cqVSpeBNS5n7C8IVInixGAoxmnlMuBnhbgrkm",
	}
	pages, err := filepath.Glob("../../../web/*.html")
	if err != nil {
		t.Fatalf("list templates: %v", err)
	}
	cdnTag := regexp.MustCompile(`<(?:script|link)\b[^>]*https://cdn\.jsdelivr\.net/[^>]*>`)
	for _, page := range pages {
		content, readErr := os.ReadFile(page) //nolint:gosec
		if readErr != nil {
			t.Fatalf("read template %s: %v", page, readErr)
		}
		for _, tag := range cdnTag.FindAllString(string(content), -1) {
			if !strings.Contains(tag, "@") {
				t.Errorf("%s contains an unpinned CDN asset: %s", filepath.Base(page), tag)
			}
			if !strings.Contains(tag, `integrity="sha384-`) || !strings.Contains(tag, `crossorigin="anonymous"`) {
				t.Errorf("%s contains a CDN asset without SRI protection: %s", filepath.Base(page), tag)
			}
			matched := false
			for url, integrity := range expectedIntegrityByURL {
				if strings.Contains(tag, url) {
					matched = true
					if !strings.Contains(tag, `integrity="`+integrity+`"`) {
						t.Errorf("%s contains an unexpected SRI value for %s", filepath.Base(page), url)
					}
					break
				}
			}
			if !matched {
				t.Errorf("%s contains an unreviewed CDN asset: %s", filepath.Base(page), tag)
			}
		}
	}
}

func TestWebAssetsDoNotReferenceCDNVendoredCopies(t *testing.T) {
	t.Parallel()
	assets, err := filepath.Glob("../../../web/*.html")
	if err != nil {
		t.Fatalf("list templates: %v", err)
	}
	assets = append(assets, "../../../web/js/waiter-service-worker.js")

	ignoredCopies := []string{
		"/css/external/bootstrap.min.css",
		"/js/external/alpine.min.js",
		"/js/external/bootstrap.bundle.min.js",
	}
	for _, asset := range assets {
		content, readErr := os.ReadFile(asset) //nolint:gosec
		if readErr != nil {
			t.Fatalf("read asset %s: %v", asset, readErr)
		}
		for _, ignoredCopy := range ignoredCopies {
			if strings.Contains(string(content), ignoredCopy) {
				t.Errorf("%s references ignored local asset %q", filepath.Base(asset), ignoredCopy)
			}
		}
	}
}

func TestFontAwesomeLicenseAndAttributionArePreserved(t *testing.T) {
	t.Parallel()

	files := map[string][]string{
		"../../../web/fontawesome/LICENSE.txt": {
			"Font Awesome Free License",
			"CC BY 4.0 License",
			"SIL OFL 1.1 License",
			"MIT License",
		},
		"../../../web/fontawesome/css/all.min.css": {
			"Font Awesome Free 5.15.4 by @fontawesome",
			"Icons: CC BY 4.0, Fonts: SIL OFL 1.1, Code: MIT License",
		},
		"../../../../THIRD_PARTY_NOTICES.md": {
			"## Font Awesome Free 5.15.4",
			"Copyright Fonticons, Inc.",
			"backend/web/fontawesome/LICENSE.txt",
		},
	}

	for path, notices := range files {
		content, err := os.ReadFile(path) //nolint:gosec
		if err != nil {
			t.Fatalf("read Font Awesome notice file %s: %v", path, err)
		}
		for _, notice := range notices {
			if !strings.Contains(string(content), notice) {
				t.Errorf("%s does not contain required notice %q", path, notice)
			}
		}
	}

	for _, font := range []string{
		"fa-regular-400.woff2",
		"fa-solid-900.woff2",
	} {
		info, err := os.Stat(filepath.Join("../../../web/fontawesome/webfonts", font))
		if err != nil {
			t.Errorf("Font Awesome runtime font %s is missing: %v", font, err)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("Font Awesome runtime font %s is empty", font)
		}
	}
}

func TestOtoLogicSoundAttributionIsPreserved(t *testing.T) {
	t.Parallel()

	content, err := os.ReadFile("../../../web/sound/chime.mp3")
	if err != nil {
		t.Fatalf("read OtoLogic sound: %v", err)
	}
	const wantHash = "717b85d29622098e92e0b74d9519d4ca8e000e28446f9e55dbb68cc827b7fa10"
	if gotHash := fmt.Sprintf("%x", sha256.Sum256(content)); gotHash != wantHash {
		t.Errorf("chime.mp3 SHA-256 = %s, want %s; update the asset and license inventory together", gotHash, wantHash)
	}

	notice, err := os.ReadFile("../../../../THIRD_PARTY_NOTICES.md")
	if err != nil {
		t.Fatalf("read third-party notices: %v", err)
	}
	for _, required := range []string{
		"効果音：OtoLogic（CC BY 4.0）",
		"Doorbell-Ding_Dong02-1",
		"https://otologic.jp/free/license",
	} {
		if !strings.Contains(string(notice), required) {
			t.Errorf("THIRD_PARTY_NOTICES.md does not contain required sound notice %q", required)
		}
	}
}

func TestVendoredFrontendLicensesMatchAssets(t *testing.T) {
	t.Parallel()

	assets := map[string]string{
		"cash.min.js":   "9a044188efdb625c5e04d1220698c099927ff16bfb434c37cd7f04dd5ee1ae1f",
		"qrcode.min.js": "41ea01fd8d9a05bb9c45e73b0b3b3e3f62dc1e44924abdff9fd39f9362f3379a",
	}
	for name, wantHash := range assets {
		content, err := os.ReadFile(filepath.Join("../../../web/js/external", name)) //nolint:gosec
		if err != nil {
			t.Fatalf("read vendored frontend asset %s: %v", name, err)
		}
		gotHash := fmt.Sprintf("%x", sha256.Sum256(content))
		if gotHash != wantHash {
			t.Errorf("%s SHA-256 = %s, want %s; update the asset and license inventory together", name, gotHash, wantHash)
		}
	}

	content, err := os.ReadFile("../../../web/js/external/LICENSES.md")
	if err != nil {
		t.Fatalf("read frontend license inventory: %v", err)
	}
	licenses := string(content)
	for _, notice := range []string{
		"## Cash 8.1.5",
		"## jquery.qrcode 1.0.0",
		"`jQuery`から`cash`へ変更",
		"## MIT License本文（英語原文）",
	} {
		if !strings.Contains(licenses, notice) {
			t.Errorf("frontend license inventory does not contain %q", notice)
		}
	}
}

func TestUnlicensedWebAssetsAreNotReferenced(t *testing.T) {
	t.Parallel()

	for _, removedPath := range []string{
		"../../../web/js/external/weights",
		"../../../web/image/allergy/egg.jpg",
		"../../../web/image/product/base_image.jpg",
		"../../../web/image/favicon.ico",
		"../../../web/image/180.png",
		"../../../web/image/192.png",
		"../../../web/image/512.png",
	} {
		if _, err := os.Stat(removedPath); !os.IsNotExist(err) {
			t.Errorf("unlicensed asset still exists: %s", removedPath)
		}
	}

	assets, err := filepath.Glob("../../../web/js/*.js")
	if err != nil {
		t.Fatalf("list application scripts: %v", err)
	}
	assets = append(assets, "../../../web/verify_twenty.html", "../../../web/js/waiter-manifest.json")
	for _, asset := range assets {
		content, readErr := os.ReadFile(asset) //nolint:gosec
		if readErr != nil {
			t.Fatalf("read asset %s: %v", asset, readErr)
		}
		for _, forbidden := range []string{
			"storage.googleapis.com/row-allergy",
			"storage.googleapis.com/row-product",
			"js/external/weights",
			"face-api.js",
		} {
			if strings.Contains(string(content), forbidden) {
				t.Errorf("%s still references removed asset %q", filepath.Base(asset), forbidden)
			}
		}
	}
}

func TestApplicationPagesUseAlpineWithoutJQuery(t *testing.T) {
	t.Parallel()
	pages := []string{
		"admin_product.html",
		"back_yard.html",
		"index.html",
		"internet.html",
		"public.html",
		"reception.html",
		"repository.html",
		"signboard.html",
		"verify_twenty.html",
		"waiter.html",
		"waiter_place.html",
	}

	for _, page := range pages {
		t.Run(page, func(t *testing.T) {
			t.Parallel()
			// pageは上記の固定リストから選択され、外部入力から渡されることはありません。
			content, err := os.ReadFile(filepath.Join("../../../web", page)) //nolint:gosec
			if err != nil {
				t.Fatalf("read template: %v", err)
			}
			template := string(content)
			if !strings.Contains(template, `x-data="restaurantOrderServerPage"`) {
				t.Error("template does not initialize the Alpine page component")
			}
			if !strings.Contains(template, "https://cdn.jsdelivr.net/npm/alpinejs@3.17.4/dist/cdn.min.js") {
				t.Error("template does not load Alpine.js")
			}
			if !strings.Contains(template, "/js/external/cash.min.js") {
				t.Error("template does not load the lightweight DOM compatibility layer")
			}
			if strings.Contains(strings.ToLower(template), "jquery") {
				t.Error("template still references jQuery")
			}
		})
	}
}

func TestIndexPageProvidesLinkQRPreview(t *testing.T) {
	t.Parallel()
	content, err := os.ReadFile("../../../web/index.html")
	if err != nil {
		t.Fatalf("read template: %v", err)
	}
	template := string(content)
	for _, expected := range []string{
		"id=\"qr-preview\"",
		"/js/external/qrcode.min.js",
		"/js/index.js",
		"{{range .PlaceGroups}}",
		"class=\"place-group\"",
		`class="fas fa-utensils"`,
		`class="fas fa-concierge-bell"`,
		`class="fas fa-mobile-alt"`,
	} {
		if !strings.Contains(template, expected) {
			t.Errorf("index template does not contain %q", expected)
		}
	}
	if strings.Contains(template, ".AssetURL") {
		t.Error("index assets must use same-origin paths")
	}
	if strings.Contains(template, "fa-solid") || !strings.Contains(template, `class="fas {{.Icon}}"`) {
		t.Error("index template must use the bundled Font Awesome 5 solid icon class")
	}
}

func TestWaiterSkipsServiceWorkerForLocalSelfSignedTLS(t *testing.T) {
	t.Parallel()
	content, err := os.ReadFile("../../../web/waiter.html")
	if err != nil {
		t.Fatalf("read waiter template: %v", err)
	}
	template := string(content)
	for _, expected := range []string{
		`'{{ .Environment }}' === 'local'`,
		`navigator.serviceWorker.getRegistrations()`,
		`registration.unregister()`,
		`catch (error)`,
	} {
		if !strings.Contains(template, expected) {
			t.Errorf("waiter template does not contain %q", expected)
		}
	}
}
