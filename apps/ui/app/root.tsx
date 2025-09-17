import {
	isRouteErrorResponse,
	Links,
	Meta,
	Outlet,
	Scripts,
	ScrollRestoration,
} from "react-router";
import "./app.css";

/**
 * External resource links for the application.
 * Includes font preconnections and stylesheets for improved performance.
 */
export const links = () => [
	{ rel: "preconnect", href: "https://fonts.googleapis.com" },
	{
		rel: "preconnect",
		href: "https://fonts.gstatic.com",
		crossOrigin: "anonymous",
	},
	{
		rel: "stylesheet",
		href: "https://fonts.googleapis.com/css2?family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&display=swap",
	},
];

/**
 * Root layout component that wraps the entire application.
 * Provides the base HTML structure including meta tags, external links,
 * and essential React Router components for navigation and script loading.
 * 
 * @param children - React components to render within the layout
 * @returns JSX element representing the complete HTML document structure
 */
export function Layout({ children }: { children: React.ReactNode }) {
	return (
		<html lang="en">
			<head>
				<meta charSet="utf-8" />
				<meta name="viewport" content="width=device-width, initial-scale=1" />
				<Meta />
				<Links />
			</head>
			<body>
				{children}
				<ScrollRestoration />
				<Scripts />
			</body>
		</html>
	);
}

/**
 * Main application component that serves as the root of the React component tree.
 * This component renders the current route using React Router's Outlet component,
 * allowing for nested routing and layout composition.
 * 
 * @returns JSX element that renders the current route's component
 */
export default function App() {
	return <Outlet />;
}
