import { RegisterForm } from '~/components/auth/register-form';

export function meta() {
  return [
    { title: 'Registreren | PWS ELO' },
    { name: 'description', content: 'Maak een nieuw PWS ELO account aan' },
  ];
}

export default function Register() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-neutral-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-112 w-full space-y-8 border border-neutral-500/20 p-6 rounded-xl">
        <div>
          <h2 className="mt-6 text-center text-3xl font-bold text-neutral-900">Account aanmaken</h2>
          <p className="mt-2 text-center text-sm text-neutral-600">
            Of{' '}
            <a href="/login" className="font-medium text-primary-600 hover:text-primary-500">
              log in met je bestaande account
            </a>
          </p>
        </div>
        <RegisterForm />
      </div>
    </div>
  );
}
