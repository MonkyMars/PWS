export interface Deadline {
	id: string;
	subjectId: string;
	ownerId: string;
	title: string;
	description: string;
	dueDate: string; // ISO date string
	createdAt: string; // ISO date string
}

export interface DeadlineFilters {
	subjectId?: string;
	ownerId?: string;
	dateFrom?: string; // ISO date string
	dateTo?: string;   // ISO date string
}