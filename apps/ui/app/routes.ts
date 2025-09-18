import { type RouteConfig, index, route } from "@react-router/dev/routes";

/**
 * Application route configuration for React Router.
 *
 * This file defines the routing structure for the PWS application using React Router v7.
 * The configuration uses file-based routing with explicit route definitions.
 *
 * Routes defined:
 * - `/` (index): Home page route that renders the home.tsx component
 * - `/login`: User authentication page
 * - `/register`: User registration page
 * - `/dashboard`: Protected dashboard page for authenticated users
 * - `/subjects/:subjectId`: Subject detail page with announcements and files
 *
 * @see {@link https://reactrouter.com/en/main/start/framework/routing} React Router routing documentation
 */
export default [
  index("routes/home.tsx"),
  route("login", "routes/login.tsx"),
  route("register", "routes/register.tsx"),
  route("dashboard", "routes/dashboard.tsx"),
  route("subjects/:subjectId", "routes/subjects.$subjectId.tsx"),
] satisfies RouteConfig;
