package main

func (queue *SongsQueue) Pop() Song {
	bot.session.RLock()
	firstSong := queue.songs[0]
	queue.songs = queue.songs[1:]
	bot.session.RUnlock()
	return firstSong
}

func (queue *SongsQueue) Push(element Song) {
	bot.session.RLock()
	songs := make([]Song, 0)
	songs = append(songs, element)

	for _, song := range queue.songs {
		songs = append(songs, song)
	}

	queue.songs = songs
	bot.session.RUnlock()
}

func (queue SongsQueue) Len() int {
	return len(queue.songs)
}

func (queue *SongsQueue) Clear() {
	queue.songs = make([]Song, 0)
}
