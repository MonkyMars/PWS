import React from 'react';

interface ButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
	variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger';
	size?: 'sm' | 'md' | 'lg';
	isLoading?: boolean;
	children: React.ReactNode;
}

const Button = React.forwardRef<HTMLButtonElement, ButtonProps>(
	({ className = '', variant = 'primary', size = 'md', isLoading = false, children, disabled, ...props }, ref) => {
		const baseClasses = 'inline-flex items-center justify-center font-medium rounded-lg transition-all duration-200 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary-500 focus-visible:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed';

		const variantClasses = {
			primary: 'bg-primary-600 text-white hover:bg-primary-700 active:bg-primary-800',
			secondary: 'bg-secondary-600 text-white hover:bg-secondary-700 active:bg-secondary-800',
			outline: 'border border-neutral-300 bg-white text-neutral-700 hover:bg-neutral-50 active:bg-neutral-100',
			ghost: 'text-neutral-700 hover:bg-neutral-100 active:bg-neutral-200',
			danger: 'bg-error-600 text-white hover:bg-error-700 active:bg-error-800',
		};

		const sizeClasses = {
			sm: 'px-3 py-1.5 text-sm',
			md: 'px-4 py-2 text-base',
			lg: 'px-6 py-3 text-lg',
		};

		const classes = `${baseClasses} ${variantClasses[variant]} ${sizeClasses[size]} ${className}`;

		return (
			<button
				ref={ref}
				className={classes}
				disabled={disabled || isLoading}
				{...props}
			>
				{isLoading && (
					<svg
						className="w-4 h-4 mr-2 animate-spin"
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
					>
						<circle
							className="opacity-25"
							cx="12"
							cy="12"
							r="10"
							stroke="currentColor"
							strokeWidth="4"
						/>
						<path
							className="opacity-75"
							fill="currentColor"
							d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
						/>
					</svg>
				)}
				{children}
			</button>
		);
	}
);

Button.displayName = 'Button';

export { Button, type ButtonProps };