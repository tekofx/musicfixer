package model

type Album struct {
	Name      string
	Songs     []Song
	MultiDisk bool
}

func (a *Album) AddSong(song Song) {
	a.Songs = append(a.Songs, song)
}
