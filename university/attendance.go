package university

type (
	Diary struct {
		Id       int64
		Date     int64
		Subject  int64
		Start    int64
		Duration int64
		Guru     int64
		Group    int64

		Marks []DiaryMarks `db:"-"`

		Created  int64
		Modified int64
		Deleted  int64
		Version  int64
	}

	DiaryMarks struct {
		Id        int64
		Date      int64
		SubjectId int64
		UserId    int64

		Val string //4,5 n
	}
)
