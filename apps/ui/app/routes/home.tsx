import { Link } from "react-router";
import {
  BookOpen,
  Users,
  FileText,
  Bell,
  ArrowRight,
  Check,
  Smartphone,
  Clock,
  Shield,
  Mail,
  Phone,
  MapPin,
} from "lucide-react";
import { Button } from "~/components/ui/button";
import { useCurrentUser } from "~/hooks";
import { Input } from "~/components";

export function meta() {
  return [
    { title: "PWS ELO - Elektronische Leeromgeving" },
    {
      name: "description",
      content:
        "De moderne elektronische leeromgeving voor middelbare scholieren. Toegang tot al je vakken, mededelingen en bestanden op één plek.",
    },
  ];
}

/**
 * Home page route component for the PWS application.
 * This component serves as the entry point for users and displays
 * the welcome screen with navigation links and application information.
 *
 * @returns JSX element rendering the Welcome component
 */
export default function Home() {
  const { data: user } = useCurrentUser();

  return (
    <div className="min-h-screen">
      {/* Hero Section */}
      <section className="relative bg-gradient-to-br from-primary-600 via-primary-700 to-primary-800 text-white overflow-hidden">
        {/* Background Pattern */}
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_30%_40%,rgba(255,255,255,0.1),transparent_70%)]" />
        <div className="absolute inset-0 bg-[radial-gradient(circle_at_70%_80%,rgba(255,255,255,0.05),transparent_50%)]" />

        {/* Floating Elements */}
        <div className="absolute top-20 left-10 w-16 h-16 bg-white/10 rounded-full blur-xl animate-float" />
        <div
          className="absolute bottom-32 right-16 w-24 h-24 bg-secondary-400/20 rounded-full blur-2xl animate-float"
          style={{ animationDelay: "1s" }}
        />
        <div
          className="absolute top-1/2 left-1/4 w-8 h-8 bg-white/15 rounded-full blur-lg animate-float"
          style={{ animationDelay: "2s" }}
        />

        <div className="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-24 md:py-32">
          <div className="text-center space-y-8">
            {/* Badge */}
            <div className="inline-flex items-center px-4 py-2 bg-white/10 backdrop-blur-sm rounded-full border border-white/20">
              <BookOpen className="h-4 w-4 mr-2" />
              <span className="text-sm font-medium">Moderne Leeromgeving</span>
            </div>

            {/* Title */}
            <div className="space-y-4">
              <h1 className="text-5xl md:text-7xl font-bold mb-4 tracking-tight">
                <span className="bg-gradient-to-r from-white to-primary-100 bg-clip-text text-transparent">
                  PWS ELO
                </span>
              </h1>
              <div className="w-24 h-1 bg-gradient-to-r from-secondary-400 to-secondary-500 rounded-full mx-auto"></div>
            </div>

            {/* Description */}
            <p className="text-xl md:text-2xl text-primary-100 max-w-4xl mx-auto leading-relaxed">
              De moderne{" "}
              <strong className="text-white">Elektronische Leeromgeving</strong>{" "}
              voor middelbare scholieren. Toegang tot al je vakken, mededelingen
              en bestanden op één plek.
            </p>

            {/* Key Features */}
            <div className="flex flex-wrap justify-center gap-6 text-sm font-medium text-primary-200">
              <div className="flex items-center space-x-2">
                <div className="w-2 h-2 bg-secondary-400 rounded-full"></div>
                <span>Overzichtelijk Dashboard</span>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-2 h-2 bg-secondary-400 rounded-full"></div>
                <span>Real-time Mededelingen</span>
              </div>
              <div className="flex items-center space-x-2">
                <div className="w-2 h-2 bg-secondary-400 rounded-full"></div>
                <span>Mobiel Vriendelijk</span>
              </div>
            </div>

            {/* Action Buttons */}
            <div className="flex flex-col sm:flex-row gap-4 justify-center pt-4">
              {user ? (
                <Link to="/dashboard">
                  <Button
                    size="lg"
                    variant="secondary"
                    className="text-lg px-8 py-3 hover:shadow-lg hover:shadow-secondary-500/25 transition-all duration-300 transform hover:scale-105"
                  >
                    Ga naar Dashboard
                    <ArrowRight className="ml-2 h-5 w-5" />
                  </Button>
                </Link>
              ) : (
                <>
                  <Link to="/register">
                    <Button
                      size="lg"
                      variant="secondary"
                      className="text-lg px-8 py-3 hover:shadow-lg hover:shadow-secondary-500/25 transition-all duration-300 transform hover:scale-105"
                    >
                      Account Aanmaken
                      <ArrowRight className="ml-2 h-5 w-5" />
                    </Button>
                  </Link>
                  <Link to="/login">
                    <Button
                      size="lg"
                      variant="outline"
                      className="text-lg px-8 py-3 border-white/50 text-black hover:text-white hover:bg-white/10 hover:border-white transition-all duration-300 backdrop-blur-sm"
                    >
                      Inloggen
                    </Button>
                  </Link>
                </>
              )}
            </div>

            {/* Trust Indicators */}
            <div className="pt-8 border-t border-white/20 mt-12">
              <p className="text-sm text-primary-200 mb-4">
                Vertrouwd door meer dan 2.500 leerlingen
              </p>
              <div className="flex justify-center items-center space-x-8 opacity-60">
                <div className="flex items-center space-x-2">
                  <Shield className="h-4 w-4" />
                  <span className="text-xs">Veilig & Privé</span>
                </div>
                <div className="flex items-center space-x-2">
                  <Clock className="h-4 w-4" />
                  <span className="text-xs">24/7 Beschikbaar</span>
                </div>
                <div className="flex items-center space-x-2">
                  <Smartphone className="h-4 w-4" />
                  <span className="text-xs">Alle Apparaten</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* What is ELO Section */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-neutral-900 mb-4">
              Wat is een ELO?
            </h2>
            <p className="text-lg text-neutral-600 max-w-8xl mx-auto">
              Een Elektronische Leeromgeving (ELO) is een digitaal platform dat
              het leren en onderwijzen ondersteunt. Het biedt een centrale plek
              waar leerlingen en docenten kunnen communiceren, materialen delen
              en leervoortgang bijhouden.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            <div className="text-center">
              <div className="w-16 h-16 bg-primary-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <BookOpen className="h-8 w-8 text-primary-600" />
              </div>
              <h3 className="text-xl font-semibold mb-2 text-neutral-900">
                Georganiseerd Leren
              </h3>
              <p className="text-neutral-600">
                Al je vakken, opdrachten en materialen op één overzichtelijke
                plek.
              </p>
            </div>

            <div className="text-center">
              <div className="w-16 h-16 bg-secondary-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Users className="h-8 w-8 text-secondary-600" />
              </div>
              <h3 className="text-xl font-semibold mb-2 text-neutral-900">
                Samenwerking
              </h3>
              <p className="text-neutral-600">
                Directe communicatie met docenten en medestudenten.
              </p>
            </div>

            <div className="text-center">
              <div className="w-16 h-16 bg-success-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <FileText className="h-8 w-8 text-success-600" />
              </div>
              <h3 className="text-xl font-semibold mb-2 text-neutral-900">
                Digitale Materialen
              </h3>
              <p className="text-neutral-600">
                Toegang tot alle studiebestanden, presentaties en documenten.
              </p>
            </div>

            <div className="text-center">
              <div className="w-16 h-16 bg-warning-100 rounded-full flex items-center justify-center mx-auto mb-4">
                <Bell className="h-8 w-8 text-warning-600" />
              </div>
              <h3 className="text-xl font-semibold mb-2 text-neutral-900">
                Real-time Updates
              </h3>
              <p className="text-neutral-600">
                Ontvang mededelingen en updates van je docenten direct.
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-20 bg-neutral-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-neutral-900 mb-4">
              Waarom PWS ELO?
            </h2>
            <p className="text-lg text-neutral-600 max-w-8xl mx-auto">
              Onze ELO is speciaal ontworpen voor middelbare scholieren en biedt
              alle tools die je nodig hebt voor succesvol leren.
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            <div className="space-y-8">
              <div className="flex items-start space-x-4">
                <div className="w-8 h-8 bg-primary-600 rounded-full flex items-center justify-center flex-shrink-0 mt-1">
                  <Check className="h-5 w-5 text-white" />
                </div>
                <div>
                  <h3 className="text-xl font-semibold mb-2 text-neutral-900">
                    Gebruiksvriendelijk Design
                  </h3>
                  <p className="text-neutral-600">
                    Intuïtieve interface die gemakkelijk te navigeren is, zelfs
                    voor nieuwe gebruikers.
                  </p>
                </div>
              </div>

              <div className="flex items-start space-x-4">
                <div className="w-8 h-8 bg-primary-600 rounded-full flex items-center justify-center flex-shrink-0 mt-1">
                  <Check className="h-5 w-5 text-white" />
                </div>
                <div>
                  <h3 className="text-xl font-semibold mb-2 text-neutral-900">
                    Mobiel Vriendelijk
                  </h3>
                  <p className="text-neutral-600">
                    Toegang tot je leeromgeving op elke device - desktop, tablet
                    of smartphone.
                  </p>
                </div>
              </div>

              <div className="flex items-start space-x-4">
                <div className="w-8 h-8 bg-primary-600 rounded-full flex items-center justify-center flex-shrink-0 mt-1">
                  <Check className="h-5 w-5 text-white" />
                </div>
                <div>
                  <h3 className="text-xl font-semibold mb-2 text-neutral-900">
                    Veilig & Betrouwbaar
                  </h3>
                  <p className="text-neutral-600">
                    Je gegevens zijn veilig met moderne
                    beveiligingstechnologieën en privacy-bescherming.
                  </p>
                </div>
              </div>

              <div className="flex items-start space-x-4">
                <div className="w-8 h-8 bg-primary-600 rounded-full flex items-center justify-center flex-shrink-0 mt-1">
                  <Check className="h-5 w-5 text-white" />
                </div>
                <div>
                  <h3 className="text-xl font-semibold mb-2 text-neutral-900">
                    24/7 Beschikbaar
                  </h3>
                  <p className="text-neutral-600">
                    Leer op jouw tempo, wanneer en waar je wilt. Altijd toegang
                    tot je materialen.
                  </p>
                </div>
              </div>
            </div>

            <div className="bg-white rounded-2xl shadow-xl p-8 border border-neutral-200">
              <div className="text-center space-y-6">
                <div className="w-20 h-20 bg-gradient-to-br from-primary-500 to-primary-600 rounded-full flex items-center justify-center mx-auto">
                  <BookOpen className="h-10 w-10 text-white" />
                </div>
                <h3 className="text-2xl font-bold text-neutral-900">
                  Klaar om te beginnen?
                </h3>
                <p className="text-neutral-600">
                  Maak je account aan en krijg direct toegang tot alle functies
                  van PWS ELO.
                </p>
                {!user && (
                  <Link to="/register">
                    <Button size="lg" className="w-full">
                      Gratis Account Aanmaken
                      <ArrowRight className="ml-2 h-5 w-5" />
                    </Button>
                  </Link>
                )}
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Contact Section */}
      <section className="py-20 bg-white">
        <div className="max-w-8xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold text-neutral-900 mb-4">
              Neem Contact Op
            </h2>
            <p className="text-lg text-neutral-600 max-w-5xl mx-auto">
              Heb je vragen over PWS ELO of wil je meer informatie? We helpen je
              graag verder!
            </p>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12">
            {/* Contact Info */}
            <div className="space-y-8">
              <div className="flex items-start space-x-4">
                <div className="w-12 h-12 bg-primary-100 rounded-lg flex items-center justify-center flex-shrink-0">
                  <Mail className="h-6 w-6 text-primary-600" />
                </div>
                <div>
                  <h3 className="text-xl font-semibold mb-2 text-neutral-900">
                    E-mail Support
                  </h3>
                  <p className="text-neutral-600 mb-2">
                    Voor technische vragen en ondersteuning
                  </p>
                  <a
                    href="mailto:support@pwsschool.nl"
                    className="text-primary-600 hover:text-primary-700 font-medium"
                  >
                    support@pwsschool.nl
                  </a>
                </div>
              </div>

              <div className="flex items-start space-x-4">
                <div className="w-12 h-12 bg-secondary-100 rounded-lg flex items-center justify-center flex-shrink-0">
                  <Phone className="h-6 w-6 text-secondary-600" />
                </div>
                <div>
                  <h3 className="text-xl font-semibold mb-2 text-neutral-900">
                    Telefoon
                  </h3>
                  <p className="text-neutral-600 mb-2">
                    Bereikbaar op werkdagen van 8:00 - 17:00
                  </p>
                  <a
                    href="tel:+31123456789"
                    className="text-primary-600 hover:text-primary-700 font-medium"
                  >
                    +31 (0)12 345 67 89
                  </a>
                </div>
              </div>

              <div className="flex items-start space-x-4">
                <div className="w-12 h-12 bg-success-100 rounded-lg flex items-center justify-center flex-shrink-0">
                  <MapPin className="h-6 w-6 text-success-600" />
                </div>
                <div>
                  <h3 className="text-xl font-semibold mb-2 text-neutral-900">
                    Bezoekadres
                  </h3>
                  <p className="text-neutral-600">
                    PWS School
                    <br />
                    Schoolstraat 123
                    <br />
                    1234 AB Schoolstad
                  </p>
                </div>
              </div>
            </div>

            {/* Contact Form */}
            <div className="bg-neutral-50 rounded-2xl p-8">
              <h3 className="text-2xl font-bold text-neutral-900 mb-6">
                Stuur ons een bericht
              </h3>
              <form className="space-y-6">
                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-neutral-700 mb-1">
                      Voornaam
                    </label>
                    <Input
                      type="text"
                      placeholder="Jan"
                      required
                      className="w-full px-3 py-2 border border-neutral-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-neutral-700 mb-1">
                      Achternaam
                    </label>
                    <Input
                      type="text"
                      placeholder="Jansen"
                      required
                      className="w-full px-3 py-2 border border-neutral-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    />
                  </div>
                </div>
                <div>
                  <label className="block text-sm font-medium text-neutral-700 mb-1">
                    E-mailadres
                  </label>
                  <Input
                    type="email"
                    placeholder="123456@chrlyceumdelft.nl"
                    required
                    className="w-full px-3 py-2 border border-neutral-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-neutral-700 mb-1">
                    Onderwerp
                  </label>
                  <select className="w-full px-3 text-neutral-700 py-2 border border-neutral-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent">
                    <option value="">Selecteer een onderwerp</option>
                    <option value="account">Account problemen</option>
                    <option value="technical">Technische ondersteuning</option>
                    <option value="general">Algemene vragen</option>
                    <option value="feedback">Feedback</option>
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-neutral-700 mb-1">
                    Bericht
                  </label>
                  <textarea
                    rows={4}
                    className="w-full px-3 py-2 text-neutral-700 border border-neutral-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                    placeholder="Beschrijf je vraag of probleem..."
                    required
                  ></textarea>
                </div>
                <Button type="submit" className="w-full">
                  <Mail className="h-4 w-4 mr-2" />
                  Bericht Versturen
                </Button>
              </form>
            </div>
          </div>
        </div>
      </section>

      {/* Stats Section */}
      <section className="py-20 bg-primary-600 text-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <h2 className="text-3xl md:text-4xl font-bold mb-4">
              PWS ELO in Cijfers
            </h2>
            <p className="text-xl text-primary-100">
              Ontdek waarom duizenden leerlingen kiezen voor onze leeromgeving.
            </p>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-4 gap-8">
            <div className="text-center">
              <div className="text-4xl md:text-5xl font-bold mb-2">2,500+</div>
              <div className="text-primary-200">Actieve Leerlingen</div>
            </div>
            <div className="text-center">
              <div className="text-4xl md:text-5xl font-bold mb-2">150+</div>
              <div className="text-primary-200">Docenten</div>
            </div>
            <div className="text-center">
              <div className="text-4xl md:text-5xl font-bold mb-2">50+</div>
              <div className="text-primary-200">Vakken</div>
            </div>
            <div className="text-center">
              <div className="text-4xl md:text-5xl font-bold mb-2">99.9%</div>
              <div className="text-primary-200">Uptime</div>
            </div>
          </div>
        </div>
      </section>
    </div>
  );
}
