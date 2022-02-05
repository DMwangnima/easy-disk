package space

type Range struct {
	Low uint64
	High uint64
}

func (rg *Range) SyncWrite(data []byte) error {

}

func (rg *Range) AsyncWrite() *Session {
    return nil
}

type Session struct {
	rg *Range
}

func (ses *Session) Transport(data []byte) {

}
