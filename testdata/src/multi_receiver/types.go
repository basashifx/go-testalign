package multi_receiver // want package:"testalign source order"

type Reader struct{}

func (r *Reader) Read() error { return nil }

func (r *Reader) Close() error { return nil }

type Writer struct{}

func (w *Writer) Write() error { return nil }

func (w *Writer) Flush() error { return nil }
