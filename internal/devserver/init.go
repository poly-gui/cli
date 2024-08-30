package devserver

type Options struct {
	ProjectPathAbs string
}

func Start(opt Options) error {
	server, err := new(opt)
	if err != nil {
		return err
	}
	return server.start()
}
