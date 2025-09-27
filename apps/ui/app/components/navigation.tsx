import { useState } from 'react';
import { Link, useLocation } from 'react-router';
import { Menu, X, BookOpen, Home, LogIn, User, LogOut } from 'lucide-react';
import { Button } from './ui/button';
import { useAuth } from '~/hooks/use-auth-context';
import { useLogout } from '~/hooks';
import { env } from '~/lib/env';

export function Navigation() {
	const [isMenuOpen, setIsMenuOpen] = useState(false);
	const location = useLocation();
	const { user, isAuthenticated } = useAuth();
	const logoutMutation = useLogout();

	const handleLogout = () => {
		logoutMutation.mutate();
		setIsMenuOpen(false);
	};

	const toggleMenu = () => setIsMenuOpen(!isMenuOpen);

	const isActivePage = (path: string) => location.pathname === path;

	const navItems = [
		{ to: '/', label: 'Home', icon: Home, requiresAuth: false },
		{
			to: '/dashboard',
			label: 'Dashboard',
			icon: BookOpen,
			requiresAuth: true,
		},
	];

	const visibleNavItems = navItems.filter((item) => !item.requiresAuth || isAuthenticated);

	return (
		<nav className="bg-white border-b border-neutral-200 sticky top-0 z-50">
			<div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
				<div className="flex justify-between items-center h-16">
					{/* Logo */}
					<Link
						to="/"
						className="flex items-center space-x-2 text-primary-600 hover:text-primary-700 transition-colors"
					>
						<BookOpen className="h-8 w-8" />
						<span className="text-xl font-bold">PWS ELO</span>
					</Link>

					{/* Desktop Navigation */}
					<div className="hidden md:flex items-center space-x-8">
						{visibleNavItems.map((item) => (
							<Link
								key={item.to}
								to={item.to}
								className={`flex items-center space-x-1 px-3 py-2 rounded-lg text-sm font-medium transition-colors ${isActivePage(item.to)
										? 'bg-primary-100 text-primary-700'
										: 'text-neutral-600 hover:text-neutral-900 hover:bg-neutral-100'
									}`}
							>
								<item.icon className="h-4 w-4" />
								<span>{item.label}</span>
							</Link>
						))}

						{/* Auth Section */}
						<div className="flex items-center space-x-3 ml-6 pl-6 border-l border-neutral-200">
							{user ? (
								<div className="flex items-center space-x-3">
									<div className="flex items-center space-x-2">
										<div className="w-8 h-8 bg-primary-100 rounded-full flex items-center justify-center">
											<User className="h-4 w-4 text-primary-600" />
										</div>
										<div className="hidden lg:block">
											<p className="text-sm font-medium text-neutral-900">{user.username}</p>
											<p className="text-xs text-neutral-500 capitalize">{user.role}</p>
										</div>
									</div>
									<Button
										variant="ghost"
										size="sm"
										onClick={handleLogout}
										isLoading={logoutMutation.isPending}
									>
										<LogOut className="h-4 w-4" />
									</Button>
								</div>
							) : (
								<div className="flex items-center space-x-2">
									{env.features.enableLogin && (
										<Link to="/login">
											<Button size="sm" variant={env.features.enableRegister ? "ghost" : "primary"}>
												<LogIn className="h-4 w-4 mr-1" />
												Inloggen
											</Button>
										</Link>
									)}
									{env.features.enableRegister && (
                    <Link to="/register">
                      <Button size="sm">Registreren</Button>
                    </Link>
                  )}
								</div>
							)}
						</div>
					</div>

					{/* Mobile menu button */}
					<Button
						variant="ghost"
						size="sm"
						className="md:hidden"
						onClick={toggleMenu}
						aria-label="Toggle menu"
					>
						{isMenuOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
					</Button>
				</div>
			</div>

			{/* Mobile Navigation */}
			{isMenuOpen && (
				<div className="md:hidden animate-slide-down">
					<div className="px-2 pt-2 pb-3 space-y-1 bg-white border-t border-neutral-200">
						{visibleNavItems.map((item) => (
							<Link
								key={item.to}
								to={item.to}
								onClick={() => setIsMenuOpen(false)}
								className={`flex items-center space-x-2 px-3 py-2 rounded-lg text-base font-medium transition-colors ${isActivePage(item.to)
										? 'bg-primary-100 text-primary-700'
										: 'text-neutral-600 hover:text-neutral-900 hover:bg-neutral-100'
									}`}
							>
								<item.icon className="h-5 w-5" />
								<span>{item.label}</span>
							</Link>
						))}

						{/* Mobile Auth Section */}
						<div className="pt-4 mt-4 border-t border-neutral-200">
							{user ? (
								<div className="space-y-3">
									<div className="flex items-center space-x-3 px-3 py-2">
										<div className="w-10 h-10 bg-primary-100 rounded-full flex items-center justify-center">
											<User className="h-5 w-5 text-primary-600" />
										</div>
										<div>
											<p className="text-base font-medium text-neutral-900">{user.username}</p>
											<p className="text-sm text-neutral-500 capitalize">{user.role}</p>
										</div>
									</div>
									<div className="px-3">
										<Button
											variant="outline"
											className="w-full justify-center"
											onClick={handleLogout}
											isLoading={logoutMutation.isPending}
										>
											<LogOut className="h-4 w-4 mr-2" />
											Uitloggen
										</Button>
									</div>
								</div>
							) : (
								<div className="space-y-2 px-3">
									{env.features.enableLogin && (
										<Link to="/login" onClick={() => setIsMenuOpen(false)}>
											<Button
												variant={env.features.enableRegister ? "outline" : undefined}
												className="w-full justify-center"
											>
												<LogIn className="h-4 w-4 mr-2" />
												Inloggen
											</Button>
										</Link>
									)}
									{env.features.enableRegister && (
										<Link to="/register" onClick={() => setIsMenuOpen(false)}>
											<Button className="w-full justify-center">Registreren</Button>
										</Link>
									)}
								</div>
							)}
						</div>
					</div>
				</div>
			)}
		</nav>
	);
}
