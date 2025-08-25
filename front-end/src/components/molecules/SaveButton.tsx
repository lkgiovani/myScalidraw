import { Button } from "@/components/ui/button";
import { LoadingSpinner } from "@/components/atoms/LoadingSpinner";
import { Save } from "lucide-react";

interface SaveButtonProps {
  onSave: () => void;
  hasUnsavedChanges: boolean;
  isSaving: boolean;
  disabled?: boolean;
}

export function SaveButton({
  onSave,
  hasUnsavedChanges,
  isSaving,
  disabled,
}: SaveButtonProps) {
  return (
    <Button
      onClick={onSave}
      disabled={!hasUnsavedChanges || isSaving || disabled}
      className="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white"
    >
      {isSaving ? (
        <>
          <LoadingSpinner size="sm" className="mr-2 border-white" />
          Salvando...
        </>
      ) : (
        <>
          <Save className="w-4 h-4 mr-2" />
          Salvar
        </>
      )}
    </Button>
  );
}
