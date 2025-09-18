import { useState } from 'react';
import { useNavigate } from 'react-router';
import { Eye, EyeOff, LogIn } from 'lucide-react';
import { Button } from '~/components/ui/button';
import { Input } from '~/components/ui/input';
import { useLogin } from '~/hooks';
import { loginSchema, type LoginFormData } from './validation';

export function LoginForm() {
	const [showPassword, setShowPassword] = useState(false);
	const [formData, setFormData] = useState<LoginFormData>({
		username: '',
		password: '',
	});
	const [errors, setErrors] = useState<Partial<Record<keyof LoginFormData, string>>>({});

	const navigate = useNavigate();
	const loginMutation = useLogin();

	const handleInputChange = (field: keyof LoginFormData) => (
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
			loginSchema.parse(formData);
			setErrors({});
			return true;
		} catch (error) {
			if (error instanceof Error && 'errors' in error) {
				const zodErrors = (error as any).errors;
				const newErrors: Partial<Record<keyof LoginFormData, string>> = {};

				zodErrors.forEach((err: any) => {
					if (err.path[0]) {
						newErrors[err.path[0] as keyof LoginFormData] = err.message;
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
			await loginMutation.mutateAsync(formData);
			navigate('/dashboard');
		} catch (error) {
			// Error is handled by the mutation
			console.error('Login failed:', error);
		}
	};

	return (
		<form className="mt-8 space-y-6" onSubmit={handleSubmit}>
			<div className="space-y-4">
				<Input
					label="Gebruikersnaam"
					type="text"
					value={formData.username}
					onChange={handleInputChange('username')}
					error={errors.username}
					placeholder="123456"
					maxLength={6}
					autoComplete="username"
				/>

				<div className="relative">
					<Input
						label="Wachtwoord"
						type={showPassword ? 'text' : 'password'}
						value={formData.password}
						onChange={handleInputChange('password')}
						error={errors.password}
						placeholder="••••••••"
						autoComplete="current-password"
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

			{loginMutation.error && (
				<div className="bg-error-50 border border-error-200 rounded-lg p-4">
					<p className="text-sm text-error-700">
						{loginMutation.error.message || 'Er is een fout opgetreden bij het inloggen.'}
					</p>
				</div>
			)}

			<div className="flex items-center justify-between">
				<div className="flex items-center">
					<input
						id="remember-me"
						name="remember-me"
						type="checkbox"
						className="h-4 w-4 text-primary-600 focus:ring-primary-500 border-neutral-300 rounded"
					/>
					<label htmlFor="remember-me" className="ml-2 block text-sm text-neutral-900">
						Onthoud mij
					</label>
				</div>

				<div className="text-sm">
					<a
						href="#"
						className="font-medium text-primary-600 hover:text-primary-500"
					>
						Wachtwoord vergeten?
					</a>
				</div>
			</div>

			<Button
				type="submit"
				className="w-full"
				isLoading={loginMutation.isPending}
				disabled={loginMutation.isPending}
			>
				<LogIn className="h-4 w-4 mr-2" />
				Inloggen
			</Button>
		</form>
	);
}