import { LoadingSpinner } from "@/components/atoms/LoadingSpinner";

interface LoadingStateProps {
  message: string;
  spinnerColor?: string;
}

export function LoadingState({ message, spinnerColor }: LoadingStateProps) {
  return (
    <div className="flex-1 flex items-center justify-center bg-explorer-main">
      <div className="text-center space-y-4">
        <LoadingSpinner
          size="lg"
          className={spinnerColor ? `border-${spinnerColor}-500` : undefined}
        />
        <p className="text-sm text-muted-foreground">{message}</p>
      </div>
    </div>
  );
}
