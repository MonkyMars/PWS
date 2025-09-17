import { Clock, Bell, FileText, User } from 'lucide-react';

export function RecentActivity() {
	// Mock data for recent activities
	const activities = [
		{
			id: 1,
			type: 'announcement',
			title: 'Nieuwe mededeling in Wiskunde',
			description: 'Huiswerk voor volgende week',
			time: '2 uur geleden',
			icon: Bell,
			color: 'text-warning-600 bg-warning-100',
		},
		{
			id: 2,
			type: 'file',
			title: 'Bestand geüpload in Natuurkunde',
			description: 'Samenvatting hoofdstuk 3.pdf',
			time: '4 uur geleden',
			icon: FileText,
			color: 'text-success-600 bg-success-100',
		},
		{
			id: 3,
			type: 'announcement',
			title: 'Nederlands - Leeslijst update',
			description: 'Nieuwe boeken toegevoegd',
			time: '1 dag geleden',
			icon: Bell,
			color: 'text-warning-600 bg-warning-100',
		},
		{
			id: 4,
			type: 'grade',
			title: 'Cijfer beschikbaar',
			description: 'Geschiedenis proefwerk: 7.5',
			time: '2 dagen geleden',
			icon: User,
			color: 'text-primary-600 bg-primary-100',
		},
	];

	const formatTimeAgo = (timeString: string) => {
		return timeString;
	};

	return (
		<div className="bg-white rounded-lg border border-neutral-200 p-6">
			<div className="flex items-center justify-between mb-4">
				<h3 className="text-lg font-semibold text-neutral-900">
					Recente Activiteiten
				</h3>
				<Clock className="h-5 w-5 text-neutral-400" />
			</div>

			<div className="space-y-4">
				{activities.map((activity) => (
					<div key={activity.id} className="flex items-start space-x-3">
						<div className={`p-2 rounded-lg ${activity.color}`}>
							<activity.icon className="h-4 w-4" />
						</div>
						<div className="flex-1 min-w-0">
							<p className="text-sm font-medium text-neutral-900 truncate">
								{activity.title}
							</p>
							<p className="text-sm text-neutral-500 truncate">
								{activity.description}
							</p>
							<p className="text-xs text-neutral-400 mt-1">
								{formatTimeAgo(activity.time)}
							</p>
						</div>
					</div>
				))}
			</div>

			<button className="w-full mt-4 text-sm text-primary-600 hover:text-primary-700 font-medium">
				Alle activiteiten bekijken
			</button>
		</div>
	);
}