import { type RouteConfig, index } from "@react-router/dev/routes";

/**
 * Application route configuration for React Router.
 *
 * This file defines the routing structure for the PWS application using React Router v7.
 * The configuration uses file-based routing with explicit route definitions.
 *
 * Routes defined:
 * - `/` (index): Home page route that renders the home.tsx component
 *
 * @see {@link https://reactrouter.com/en/main/start/framework/routing} React Router routing documentation
 */
export default [index("routes/home.tsx")] satisfies RouteConfig;
