interface SpinnerProps {
  size?: number;
  className?: string;
}

export function Spinner({ size, className = '' }: SpinnerProps) {
  if (size) {
    return (
      <div
        style={{ width: size, height: size }}
        className={`border-2 border-current border-t-transparent rounded-full animate-spin ${className}`}
      />
    );
  }

  return (
    <div className="flex justify-center items-center py-12">
      <div className={`w-8 h-8 border-4 border-orange-600 border-t-transparent rounded-full animate-spin ${className}`} />
    </div>
  );
}
