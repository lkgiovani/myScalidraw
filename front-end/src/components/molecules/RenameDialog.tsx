import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useFileStore, FileNode } from "@/stores/fileStore";
import { useRenameFile } from "@/hooks/useFileApi";
import { FileText } from "lucide-react";

interface RenameDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  file: FileNode;
}

export function RenameDialog({ open, onOpenChange, file }: RenameDialogProps) {
  const [name, setName] = useState("");
  const renameItem = useFileStore((state) => state.renameItem);
  const renameFileMutation = useRenameFile();

  const handleRename = async () => {
    if (!name.trim()) return;

    const finalName =
      file.type === "file" && !name.trim().endsWith(".excalidraw")
        ? `${name.trim()}.excalidraw`
        : name.trim();

    try {
      await renameFileMutation.mutateAsync({ id: file.id, name: finalName });

      renameItem(file.id, finalName);
      setName("");
      onOpenChange(false);
    } catch (error) {
      console.error("Erro ao renomear arquivo:", error);

      renameItem(file.id, finalName);
      setName("");
      onOpenChange(false);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleRename();
    }
  };

  const handleOpenChange = (isOpen: boolean) => {
    if (isOpen) {
      const baseName =
        file.type === "file" && file.name.endsWith(".excalidraw")
          ? file.name.replace(".excalidraw", "")
          : file.name;
      setName(baseName);
    } else {
      setName("");
    }
    onOpenChange(isOpen);
  };

  return (
    <Dialog open={open} onOpenChange={handleOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <FileText className="w-5 h-5" />
            Renomear {file.type === "file" ? "Arquivo" : "Pasta"}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4 py-4">
          <div className="space-y-2">
            <Label htmlFor="rename">Novo nome</Label>
            <Input
              id="rename"
              placeholder={file.type === "file" ? "Meu desenho" : "Nova pasta"}
              value={name}
              onChange={(e) => setName(e.target.value)}
              onKeyDown={handleKeyDown}
              autoFocus
            />
            {file.type === "file" && (
              <p className="text-xs text-muted-foreground">
                A extensão .excalidraw será mantida automaticamente
              </p>
            )}
          </div>
        </div>

        <div className="flex justify-end gap-2">
          <Button variant="outline" onClick={() => handleOpenChange(false)}>
            Cancelar
          </Button>
          <Button
            onClick={handleRename}
            disabled={!name.trim()}
            className="bg-primary hover:bg-primary-hover"
          >
            Renomear
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
