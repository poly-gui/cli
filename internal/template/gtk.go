package template

import "path/filepath"

// GTKSourceFiles is a list of templates of files that will be in the "gtk" directory in a Poly project.
var GTKSourceFiles = []templateFile{GTKCMakeLists, GTKCxxMainFile, GTKLaunchScript, GTKRPMSpec, GTKBuildScript}

var GTKCMakeLists = templateFile{
	FilePathRel:  "CMakeLists.txt",
	TemplateName: "GTKCMakeLists",
	Template: `cmake_minimum_required(VERSION 3.25.2)

project({{.AppName}} LANGUAGES C CXX)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED True)

set(RESOURCE_DIR ../build/res)
set(HAS_RESOURCES EXISTS ${GRESOURCES_DIR})

# LIB_INSTALL_DIR is set in prod build
# in debug build, we use CMAKE_CURRENT_BINARY_DIR as the LIB_INSTALL_DIR
# to simulate the install directory of the program which is usually /usr/lib64/${prog-name}
IF (NOT DEFINED LIB_INSTALL_DIR)
    set(LIB_INSTALL_DIR ${CMAKE_CURRENT_BINARY_DIR})
    set(APP_DIR ${LIB_INSTALL_DIR})

    message("Copying portable layer binary...")
    file(COPY_FILE ../build/bundle ${CMAKE_CURRENT_BINARY_DIR}/bundle)
ELSE()
    set(APP_DIR ${LIB_INSTALL_DIR}/{{.AppName}})
ENDIF ()

# setup gtkmm
find_package(PkgConfig)
pkg_check_modules(GTKMM gtkmm-4.0)
link_directories(${GTKMM_LIBRARY_DIRS})
include_directories(${GTKMM_INCLUDE_DIRS})

# setup poly
add_subdirectory(lib/gtk-poly)

# compile resources
set(GRESOURCES_XML gresources.xml)
set(GRESOURCES_OUTPUT_SRC src/res.c)
set(GRESOURCES_OUTPUT_HEADER src/res.h)
set(GRESOURCES_DIR ${CMAKE_CURRENT_SOURCE_DIR}/../res)
get_filename_component(GRESOURCES_DIR ${GRESOURCES_DIR} ABSOLUTE)

find_program(GLIB_COMPILE_RESOURCES NAMES glib-compile-resources REQUIRED)

IF (${HAS_RESOURCES})
    execute_process(
            WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
            COMMAND ${GLIB_COMPILE_RESOURCES} --generate-dependencies ${GRESOURCES_XML}
            OUTPUT_VARIABLE GRESOURCES_DEPS_REL
    )
    string(REPLACE "\n" " " GRESOURCES_DEPS_REL ${GRESOURCES_DEPS_REL})
    separate_arguments(GRESOURCES_DEPS_REL)
    set(GRESOURCES_DEPS "")
    foreach (DEP ${GRESOURCES_DEPS_REL})
        list(APPEND GRESOURCES_DEPS ${GRESOURCES_DIR}/${DEP})
    endforeach ()
    foreach (DEP ${GRESOURCES_DEPS})
        message("Resource detected: ${DEP}")
    endforeach ()

    add_custom_command(
            OUTPUT ${CMAKE_CURRENT_SOURCE_DIR}/${GRESOURCES_OUTPUT_SRC}
            WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
            COMMAND ${GLIB_COMPILE_RESOURCES}
            ARGS
            --generate-source
            --manual-register
            --sourcedir=${GRESOURCES_DIR}
            --target=${CMAKE_CURRENT_SOURCE_DIR}/${GRESOURCES_OUTPUT_SRC}
            ${GRESOURCES_XML}
            VERBATIM
            MAIN_DEPENDENCY ${GRESOURCES_XML}
            DEPENDS ${GRESOURCES_DEPS}
    )
    add_custom_command(
            OUTPUT ${CMAKE_CURRENT_SOURCE_DIR}/${GRESOURCES_OUTPUT_HEADER}
            WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
            COMMAND ${GLIB_COMPILE_RESOURCES}
            ARGS
            --generate-header
            --manual-register
            --sourcedir=${GRESOURCES_DIR}
            --target=${CMAKE_CURRENT_SOURCE_DIR}/${GRESOURCES_OUTPUT_HEADER}
            ${GRESOURCES_XML}
            VERBATIM
            MAIN_DEPENDENCY ${GRESOURCES_XML}
            DEPENDS ${GRESOURCES_DEPS}
    )

    add_custom_target(
            {{.AppName}}-resources
            DEPENDS
            ${CMAKE_CURRENT_SOURCE_DIR}/${GRESOURCES_OUTPUT_HEADER}
            ${CMAKE_CURRENT_SOURCE_DIR}/${GRESOURCES_OUTPUT_SRC})

    set_source_files_properties(
            ${CMAKE_CURRENT_SOURCE_DIR}/${GRESOURCES_OUTPUT_HEADER}
            ${CMAKE_CURRENT_SOURCE_DIR}/${GRESOURCES_OUTPUT_SRC}
            PROPERTIES GENERATED TRUE)
ENDIF ()

add_executable({{.AppName}}
        src/main.cxx)

target_link_libraries({{.AppName}} gtkpoly ${GTKMM_LIBRARIES})
IF (${HAS_RESOURCES})
    add_dependencies({{.AppName}} {{.AppName}}-resources)
ENDIF ()
target_compile_definitions({{.AppName}} PRIVATE APP_DIR=\"${APP_DIR}\")

install(TARGETS {{.AppName}} RUNTIME DESTINATION ${APP_DIR})
IF (EXISTS packaging/bundle)
    install(PROGRAMS packaging/bundle DESTINATION ${APP_DIR})
ELSE ()
    install(PROGRAMS ../build/bundle DESTINATION ${APP_DIR})
ENDIF ()
install(PROGRAMS packaging/launch.sh
        RENAME {{.AppName}}
        DESTINATION bin)
`,
}

var GTKLaunchScript = templateFile{
	FilePathRel: filepath.Join("packaging", "launch.sh"),
	Template: `#!/bin/bash

ARCH=$(uname -m)
case $ARCH in
	x86_64 | s390x | sparc64)
		LIB_DIR="/usr/lib64"
		SECONDARY_LIB_DIR="/usr/lib"
		;;
	* )
		LIB_DIR="/usr/lib"
		SECONDARY_LIB_DIR="/usr/lib64"
		;;
esac

BIN_NAME="{{.AppName}}"

if [ ! -r "$LIB_DIR"/{{.AppName}}/$BIN_NAME ]; then
    if [ ! -r $SECONDARY_LIB_DIR/{{.AppName}}/$BIN_NAME ]; then
	  echo "Error: $LIB_DIR/{{.AppName}}/$BIN_NAME not found"
	  if [ -d $SECONDARY_LIB_DIR ]; then
	    echo "       $SECONDARY_LIB_DIR/{{.AppName}}/$BIN_NAME not found"
	  fi
	  exit 1
    fi
    LIB_DIR="$SECONDARY_LIB_DIR"
fi

PROG_DIR="$LIB_DIR/{{.AppName}}"
PROG_BIN="$PROG_DIR/$BIN_NAME"

exec $PROG_BIN "$@"
`,
	TemplateName: "GTKLaunchScript",
}

var GTKRPMSpec = templateFile{
	FilePathRel: filepath.Join("packaging", "rpm", "app.spec"),
	Template: `Name:           {{.AppName}}
Version:        1.0.0
Release:        1%{?dist}
Summary:        An application made with Poly.

License:        GPLv3+
URL:            https://www.example.com/%{name}
Source0:        https://www.example.com/%{name}/releases/%{name}-%{version}.tar.gz

BuildRequires:  cmake ninja-build

%description
An application made with Poly.

%prep
%setup -q

%build
%cmake
%cmake_build

%install
%cmake_install


%files
%{_bindir}/%{name}
%dir %{_libdir}/%{name}

%changelog
- Initial
`,
	TemplateName: "GTKRPMSpec",
}

var GTKCxxMainFile = templateFile{
	FilePathRel: filepath.Join("src", "main.cxx"),
	Template: `#include <gtkpoly/application.hxx>

int main()
{
    const Poly::ApplicationConfig config{
        .application_id = "{{.PackageName}}.{{.AppName}}",
        .app_dir_path = APP_DIR,
        .flags = Gio::Application::Flags::NONE,
    };
    const auto app = Poly::Application::create(config);

    return app->start();
}
`,
	TemplateName: "GTKCxxMainFile",
}

var GTKBuildScript = templateFile{
	FilePathRel: filepath.Join("build.sh"),
	Template: `#!/bin/bash

set -eu
pushd "$(dirname $0)" > /dev/null

for arg in "$@"; do declare $arg=1; done

prefix="/usr/local/libexec"
if [ -v out ]; then prefix="${out}"; fi

if [ ! -v release ]; then
	debug=1
	app_dir="$(pwd)/build"
else
	app_dir="${out}/{{.AppName}}"
fi

# when building in nix, the dependencies will have already been built and installed
# so there is no need to fetch the submodules
if [ ! -v nix ]; then
	echo "fetching submodule dependencies..."
	git submodule update --init --recursive
fi

if [ ! -v nix ]; then
	./lib/gtk-poly/build.sh
	./lib/cxx-nanopack/build.sh
fi

link_opts="-lgtkpoly -lnanopack"
include_opts=""
if [ ! -v nix ]; then 
	link_opts="-L../lib/gtk-poly/build -L../lib/cxx-nanopack/build"
	include_opts="-I../lib/gtk-poly/include -I../lib/cxx-nanopack/include"
fi

gtkmm_flags="$(pkg-config --cflags --libs gtkmm-4.0)"
common_opts="--std=c++20 -Wall -DAPP_DIR=\"${app_dir}\" ${gtkmm_flags}"
if [ ! -v nix ]; then common_opts="${common_opts} ${include_opts} ${link_opts}"; fi
debug_opts="--debug --optimize -DDEBUG ${common_opts}"
release_opts="--optimize -DDEBUG=0 ${common_opts}"
compiler="${CC:-g++}"
ar="${CC:-ar}"
src_files=(
	src/main.cxx
)

if [ -v debug ]; then compile="$compiler ${debug_opts}"; fi
if [ -v release ]; then compile="$compiler ${release_opts}"; fi

echo "compiling using:"
echo $compile

mkdir -p build
pushd build

all_src=""
for p in "${src_files[@]}"; do
	all_src+=" ../${p}"
done

$compile -o {{.AppName}} ${all_src} -lgtkpoly -lnanopack
cp ../../build/bundle $app_dir

popd > /dev/null
popd > /dev/null
`,
	TemplateName: "GTKBuildScript",
}
