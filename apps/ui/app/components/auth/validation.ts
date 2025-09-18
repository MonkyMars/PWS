import { z } from "zod";

// Login validation schema
export const loginSchema = z.object({
  username: z
    .string()
    .min(1, "Gebruikersnaam is verplicht")
    .regex(/^\d{6}$/, "Gebruikersnaam moet uit 6 cijfers bestaan"),
  password: z
    .string()
    .min(1, "Wachtwoord is verplicht")
    .min(8, "Wachtwoord moet minimaal 8 karakters lang zijn"),
});

// Register validation schema
export const registerSchema = z.object({
  username: z
    .string()
    .min(1, "Gebruikersnaam is verplicht")
    .regex(/^\d{6}$/, "Gebruikersnaam moet uit 6 cijfers bestaan"),
  email: z
    .string()
    .min(1, "E-mailadres is verplicht")
    .email("Ongeldig e-mailadres"),
  password: z
    .string()
    .min(1, "Wachtwoord is verplicht")
    .min(8, "Wachtwoord moet minimaal 8 karakters lang zijn")
    .regex(
      /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)/,
      "Wachtwoord moet minimaal één kleine letter, één hoofdletter en één cijfer bevatten"
    ),
  firstName: z
    .string()
    .min(1, "Voornaam is verplicht")
    .min(2, "Voornaam moet minimaal 2 karakters lang zijn"),
  lastName: z
    .string()
    .min(1, "Achternaam is verplicht")
    .min(2, "Achternaam moet minimaal 2 karakters lang zijn"),
});

export type LoginFormData = z.infer<typeof loginSchema>;
export type RegisterFormData = z.infer<typeof registerSchema>;
