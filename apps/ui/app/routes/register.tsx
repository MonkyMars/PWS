import { RegisterForm } from '~/components/auth/register-form';
import { FeatureDisabled } from '~/components/feature-disabled';
import { useAuth } from '~/hooks';
import { Navigate } from 'react-router';
import { env } from '~/lib/env';

export function meta() {
	return [
		{ title: 'Registreren | PWS ELO' },
		{ name: 'description', content: 'Maak een nieuw PWS ELO account aan' },
	];
}

export default function Register() {
	const { isAuthenticated, isLoading } = useAuth();

	// Check if registration feature is disabled
	if (!env.features.enableRegister) {
		return (
			<FeatureDisabled
				featureName="Registratie"
				description="De registratie functionaliteit is momenteel uitgeschakeld door de beheerder."
			/>
		);
	}

	if (isLoading) {
		return (
			<div className="min-h-screen flex items-center justify-center">
				<div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
			</div>
		);
	}

	if (isAuthenticated) {
		return <Navigate to="/dashboard" replace />;
	}

	return (
		<div className="min-h-screen flex items-center justify-center bg-neutral-50 py-12 px-4 sm:px-6 lg:px-8">
			<div className="max-w-4xl w-full space-y-8">
				<div>
					<h2 className="mt-6 text-center text-3xl font-bold text-neutral-900">Account aanmaken</h2>
					<p className="mt-2 text-center text-sm text-neutral-600">
						Of{' '}
						{env.features.enableLogin ? (
							<a href="/login" className="font-medium text-primary-600 hover:text-primary-500">
								log in met je bestaande account
							</a>
						) : (
							<span className="text-neutral-400">inloggen is uitgeschakeld</span>
						)}
					</p>
				</div>
				<RegisterForm />
			</div>
		</div>
	);
}
