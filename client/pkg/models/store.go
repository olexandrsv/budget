package models

type Store[T Object[T]] interface {
	ForEach(func(T) error) error
	Get(int) T
	PutAll([]T)
	Put(int, T)
	Delete(int)
	Len() int
	Add(T)
}

type store[T Object[T]] struct {
	max     int
	data    map[int]T
	indexes []int
}

func NewStore[T Object[T]]() Store[T] {
	return &store[T]{
		data: make(map[int]T),
	}
}

func (s store[T]) Get(index int) T {
	return s.data[index]
}

func (s *store[T]) PutAll(items []T) {
	indexes := make([]int, 0, len(items))
	for _, item := range items {
		if s.max < item.Index() {
			s.max = item.Index()
		}
		s.data[item.Index()] = item
		indexes = append(indexes, item.Index())
	}
	s.indexes = indexes
}

func (s *store[T]) Put(index int, item T) {
	if s.max < item.Index() {
		s.max = item.Index()
	}
	if _, ok := s.data[index]; !ok {
		s.indexes = append(s.indexes, item.Index())
	}
	s.data[index] = item
}

func (s *store[T]) Delete(index int) {
	delete(s.data, index)
}

func (s *store[T]) ForEach(fn func(T) error) error {
	for _, idx := range s.indexes {
		err := fn(s.data[idx])
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *store[T]) Len() int {
	if s.data == nil {
		return 0
	}
	return len(s.data)
}

func (s *store[T]) Add(t T) {
	s.max++
	t.SetIndex(s.max)
	s.data[s.max] = t
	s.indexes = append(s.indexes, s.max)
}
