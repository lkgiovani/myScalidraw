import { Check } from "lucide-react";

interface StatusIndicatorsProps {
  isSaving: boolean;
  hasUnsavedChanges: boolean;
  selectedFileId: string | null;
  isPending: boolean;
}

export function StatusIndicators({
  isSaving,
  hasUnsavedChanges,
  selectedFileId,
  isPending,
}: StatusIndicatorsProps) {
  if (isSaving) {
    return (
      <div className="absolute top-16 right-4 z-10 bg-blue-500 text-white px-3 py-1 rounded-md text-sm font-medium">
        Auto-salvando...
      </div>
    );
  }

  if (hasUnsavedChanges && !isSaving) {
    return (
      <div className="absolute top-16 right-4 z-10 bg-yellow-500 text-yellow-900 px-3 py-1 rounded-md text-sm font-medium">
        Alterações não salvas
      </div>
    );
  }

  if (!hasUnsavedChanges && !isSaving && selectedFileId && !isPending) {
    return (
      <div className="absolute top-16 right-4 z-10 bg-green-500 text-white px-3 py-1 rounded-md text-sm font-medium flex items-center gap-1">
        <Check className="w-3 h-3" />
        Salvo
      </div>
    );
  }

  return null;
}
