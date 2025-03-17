package webhook

type Webhook struct {
}

func New(port int) *Webhook {

}

func (w *Webhook) MustRun() {
	if err := w.Run(); err != nil {
		panic(err)
	}
}

func (w *Webhook) Run() error {

	return nil // temp
}
