import { useState } from 'react';
import { useNavigate } from 'react-router';
import { Eye, EyeOff, UserPlus } from 'lucide-react';
import { Button } from '~/components/ui/button';
import { Input } from '~/components/ui/input';
import { useRegister } from '~/hooks';
import { registerSchema, type RegisterFormData } from './validation';

export function RegisterForm() {
	const [showPassword, setShowPassword] = useState(false);
	const [formData, setFormData] = useState<RegisterFormData>({
		username: '',
		email: '',
		password: '',
		firstName: '',
		lastName: '',
	});
	const [errors, setErrors] = useState<Partial<Record<keyof RegisterFormData, string>>>({});

	const navigate = useNavigate();
	const registerMutation = useRegister();

	const handleInputChange = (field: keyof RegisterFormData) => (
		e: React.ChangeEvent<HTMLInputElement>
	) => {
		const value = e.target.value;
		setFormData(prev => ({ ...prev, [field]: value }));

		// Clear error when user starts typing
		if (errors[field]) {
			setErrors(prev => ({ ...prev, [field]: undefined }));
		}
	};

	const validateForm = () => {
		try {
			registerSchema.parse(formData);
			setErrors({});
			return true;
		} catch (error) {
			if (error instanceof Error && 'errors' in error) {
				const zodErrors = (error as any).errors;
				const newErrors: Partial<Record<keyof RegisterFormData, string>> = {};

				zodErrors.forEach((err: any) => {
					if (err.path[0]) {
						newErrors[err.path[0] as keyof RegisterFormData] = err.message;
					}
				});

				setErrors(newErrors);
			}
			return false;
		}
	};

	const handleSubmit = async (e: React.FormEvent) => {
		e.preventDefault();

		if (!validateForm()) {
			return;
		}

		try {
			await registerMutation.mutateAsync(formData);
			navigate('/dashboard');
		} catch (error) {
			// Error is handled by the mutation
			console.error('Registration failed:', error);
		}
	};

	return (
		<form className="mt-8 space-y-6" onSubmit={handleSubmit}>
			<div className="space-y-4">
				<div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
					<Input
						label="Voornaam"
						type="text"
						value={formData.firstName}
						onChange={handleInputChange('firstName')}
						error={errors.firstName}
						placeholder="Jan"
						autoComplete="given-name"
					/>

					<Input
						label="Achternaam"
						type="text"
						value={formData.lastName}
						onChange={handleInputChange('lastName')}
						error={errors.lastName}
						placeholder="Jansen"
						autoComplete="family-name"
					/>
				</div>

				<Input
					label="Gebruikersnaam"
					type="text"
					value={formData.username}
					onChange={handleInputChange('username')}
					error={errors.username}
					placeholder="123456"
					helperText="Je leerlingnummer (6 cijfers)"
					maxLength={6}
					autoComplete="username"
				/>

				<Input
					label="E-mailadres"
					type="email"
					value={formData.email}
					onChange={handleInputChange('email')}
					error={errors.email}
					placeholder="jan.jansen@student.pwsschool.nl"
					autoComplete="email"
				/>

				<div className="relative">
					<Input
						label="Wachtwoord"
						type={showPassword ? 'text' : 'password'}
						value={formData.password}
						onChange={handleInputChange('password')}
						error={errors.password}
						placeholder="••••••••"
						helperText="Minimaal 8 karakters, met hoofdletter, kleine letter en cijfer"
						autoComplete="new-password"
					/>
					<button
						type="button"
						className="absolute right-3 top-9 text-neutral-400 hover:text-neutral-600 transition-colors"
						onClick={() => setShowPassword(!showPassword)}
						aria-label={showPassword ? 'Wachtwoord verbergen' : 'Wachtwoord tonen'}
					>
						{showPassword ? <EyeOff className="h-5 w-5" /> : <Eye className="h-5 w-5" />}
					</button>
				</div>
			</div>

			{registerMutation.error && (
				<div className="bg-error-50 border border-error-200 rounded-lg p-4">
					<p className="text-sm text-error-700">
						{registerMutation.error.message || 'Er is een fout opgetreden bij het registreren.'}
					</p>
				</div>
			)}

			<div className="flex items-center">
				<input
					id="accept-terms"
					name="accept-terms"
					type="checkbox"
					required
					className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-neutral-300 rounded"
				/>
				<label htmlFor="accept-terms" className="ml-2 block text-sm text-neutral-900">
					Ik ga akkoord met de{' '}
					<a href="#" className="text-primary-600 hover:text-primary-500">
						gebruiksvoorwaarden
					</a>{' '}
					en het{' '}
					<a href="#" className="text-primary-600 hover:text-primary-500">
						privacybeleid
					</a>
				</label>
			</div>

			<Button
				type="submit"
				className="w-full"
				isLoading={registerMutation.isPending}
				disabled={registerMutation.isPending}
			>
				<UserPlus className="h-4 w-4 mr-2" />
				Account aanmaken
			</Button>
		</form>
	);
}