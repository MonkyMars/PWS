import React from 'react';
import { useNavigate } from 'react-router';
import { Button } from '~/components/ui/button';
import { AlertCircle } from 'lucide-react';

export default function NotFoundPage() {
  const navigate = useNavigate();

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-neutral-50 px-4">
      <div className="bg-white rounded-lg border border-neutral-200 p-10 shadow-md flex flex-col items-center">
        <AlertCircle className="h-16 w-16 text-red-400 mb-4" />
        <h1 className="text-3xl font-bold text-neutral-900 mb-2">404 - Pagina niet gevonden</h1>
        <p className="text-neutral-600 mb-6 text-center">
          De pagina die je zoekt bestaat niet of is verplaatst.
          <br />
          Controleer het adres of ga terug naar het dashboard.
        </p>
        <Button variant="primary" onClick={() => navigate('/dashboard')}>
          Terug naar dashboard
        </Button>
      </div>
    </div>
  );
}
