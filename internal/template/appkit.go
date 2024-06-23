package template

import "path/filepath"

// AppKitSourceFiles is a list of templates of files that will be in the "macOS" directory in a Poly project.
var AppKitSourceFiles = []templateFile{XcodeGenSpec, SwiftMainFile, AppDelegate}

var XcodeGenSpec = templateFile{
	FilePathRel: "project.yml",
	Template: `name: {{.AppName}}
options:
  bundleIdPrefix: {{.PackageName}}
packages:
  PolyNative:
    {{- if .DebugWorkspacePath}}
    path: {{.DebugWorkspacePath}}/PolyNativeSwift
    {{- else}}
    url: https://github.com/poly-gui/swift-poly-native
    branch: main
    {{- end}}
settings:
  GENERATE_INFOPLIST_FILE: YES
targets:
  {{.AppName}}:
    type: application
    platform: macOS
    deploymentTarget: "10.15"
    sources: [{{.AppName}}]
    dependencies:
      - package: PolyNative
    postCompileScripts:
      - script: |
          mkdir -p "${TARGET_BUILD_DIR}/${UNLOCALIZED_RESOURCES_FOLDER_PATH}"
		  {{- if .DebugWorkspacePath}}
		  ln -s "${SRCROOT}/../build/bundle" "${TARGET_BUILD_DIR}/${UNLOCALIZED_RESOURCES_FOLDER_PATH}/bundle"
		  {{- else}}
          cp "${SRCROOT}/../build/bundle" "${TARGET_BUILD_DIR}/${UNLOCALIZED_RESOURCES_FOLDER_PATH}/bundle"
		  {{- end}}
`,
	TemplateName: "XcodeGenSpec",
}

var SwiftMainFile = templateFile{
	FilePathRel: filepath.Join("_APP_NAME_", "main.swift"),
	Template: `import AppKit

let application = NSApplication.shared
let delegate = AppDelegate()
application.delegate = delegate

_ = NSApplicationMain(CommandLine.argc, CommandLine.unsafeArgv)
`,
	TemplateName: "SwiftMainFile",
}

var AppDelegate = templateFile{
	FilePathRel: filepath.Join("_APP_NAME_", "AppDelegate.swift"),
	Template: `import PolyNative

class AppDelegate: PolyApplicationDelegate {}
`,
	TemplateName: "SwiftAppDelegate",
}
