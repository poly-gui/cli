package template

const XcodeGenSpec = `name: {{.AppName}}
options:
  bundleIdPrefix: {{.PackageName}}
packages:
  PolyNative:
    url: https://github.com/poly-gui/swift-poly-native
    branch: main
settings:
  GENERATE_INFOPLIST_FILE: YES
targets:
  {{.AppName}}:
    type: application
    platform: macOS
    deploymentTarget: "10.13"
    sources: [{{.AppName}}]
    dependencies:
      - package: PolyNative
    postCompileScripts:
      - script: |
          mkdir -p "${TARGET_BUILD_DIR}/${UNLOCALIZED_RESOURCES_FOLDER_PATH}"
          cp "${SRCROOT}/../build/bundle" "${TARGET_BUILD_DIR}/${UNLOCALIZED_RESOURCES_FOLDER_PATH}/bundle"
`

const SwiftMainFile = `import AppKit

let application = NSApplication.shared
let delegate = AppDelegate()
application.delegate = delegate

_ = NSApplicationMain(CommandLine.argc, CommandLine.unsafeArgv)
`

const AppDelegate = `import PolyNative

class AppDelegate: PolyApplicationDelegate {}
`
