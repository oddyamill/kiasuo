export interface Student {
	studentNameAcronym: string
}

export interface Marks {
	studyPeriod: StudyPeriod
	lessons: Lesson[]
	showPasses: boolean
	showEmptyLessons: boolean
	lastMarksSeenAt: number
}

interface StudyPeriod {
	text: string
}

export interface Lesson {
	subject: string
	slots: Slot[]
	averages: Averages
}

export interface Slot {
	lesson_date: string
	mark: Mark
	text: string;
	updated_at: string
}

export interface Mark {
	value: string
}

interface Averages {
	for_student: [string, number] | null
	for_class: [string, number] | null
}
