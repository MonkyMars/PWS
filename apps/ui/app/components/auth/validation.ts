import { z } from 'zod';

// Login validation schema
export const loginSchema = z.object({
  email: z.string().min(1, 'E-mailadres is verplicht').email('Ongeldig e-mailadres'),
  password: z
    .string()
    .min(1, 'Wachtwoord is verplicht')
    .min(8, 'Wachtwoord moet minimaal 8 karakters lang zijn'),
});

// Register validation schema
export const registerSchema = z.object({
  username: z
    .string()
    .min(1, 'Gebruikersnaam is verplicht')
    .min(3, 'Gebruikersnaam moet minimaal 3 karakters lang zijn'),
  email: z.string().min(1, 'E-mailadres is verplicht').email('Ongeldig e-mailadres'),
  password: z
    .string()
    .min(1, 'Wachtwoord is verplicht')
    .min(8, 'Wachtwoord moet minimaal 8 karakters lang zijn')
    .regex(
      /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/,
      'Wachtwoord moet minimaal één kleine letter, één hoofdletter en één cijfer bevatten'
    ),
});

export type LoginFormData = z.infer<typeof loginSchema>;
export type RegisterFormData = z.infer<typeof registerSchema>;
