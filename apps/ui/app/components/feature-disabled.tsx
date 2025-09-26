import { Link } from 'react-router';
import { AlertTriangle, Home } from 'lucide-react';
import { Button } from './ui/button';

interface FeatureDisabledProps {
	featureName: string;
	description?: string;
}

/**
 * Component to display when a feature is disabled via environment variables.
 */
export function FeatureDisabled({ featureName, description }: FeatureDisabledProps) {
	return (
		<div className="min-h-screen flex items-center justify-center bg-neutral-50 py-12 px-4 sm:px-6 lg:px-8">
			<div className="max-w-96 w-full space-y-8 text-center">
				<div>
					<div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-yellow-100">
						<AlertTriangle className="h-6 w-6 text-yellow-600" aria-hidden="true" />
					</div>
					<h2 className="mt-6 text-center text-3xl font-bold text-neutral-900">
						Feature Uitgeschakeld
					</h2>
					<p className="mt-2 text-center text-sm text-neutral-600">
						{description || `De ${featureName.toLowerCase()} functionaliteit is momenteel uitgeschakeld.`}
					</p>
				</div>
				<div>
					<Link to="/">
						<Button className="w-full flex justify-center items-center">
							<Home className="h-4 w-4 mr-2" />
							Terug naar Home
						</Button>
					</Link>
				</div>
			</div>
		</div>
	);
}