package cli

import (
	"errors"
	"log"
	"os"
	"path"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/fsnotify/fsnotify"
)

func Run() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if _, err := os.Stat(path.Join(cwd, "macOS")); os.IsNotExist(err) {
		return errors.New("Cannot find macOS/ directory in the current directory! Make sure you are in a Poly project before running this command.")
	}
	if _, err := os.Stat(path.Join(cwd, "app")); os.IsNotExist(err) {
		return errors.New("Cannot find app/ directory in the current directory! Make sure you are in a Poly project before running this command.")
	}

	ctx, err := api.Context(api.BuildOptions{
		EntryPoints:       []string{"src/main.ts"},
		Outfile:           "build/out.js",
		AbsWorkingDir:     path.Join(cwd, "app"),
		Bundle:            true,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Format:            api.FormatCommonJS,
		Platform:          api.PlatformNode,
	})
	if err != nil {
		return err
	}

	err = ctx.Watch(api.WatchOptions{})
	if err != nil {
		return err
	}

	err = startFSWatcher(path.Join(cwd, "app", "build", "out.js"))
	if err != nil {
		return err
	}

	<-make(chan struct{})

	return nil
}

func startFSWatcher(jsPath string) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = watcher.Add(jsPath)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error: ", err)
			}
		}
	}()

	return nil
}
