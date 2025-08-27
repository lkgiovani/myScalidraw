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
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { useFileStore } from "@/stores/fileStore";
import { FileText, Folder } from "lucide-react";

import { toast } from "sonner";

interface CreateItemDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  parentId?: string;
}

export function CreateItemDialog({
  open,
  onOpenChange,
  parentId,
}: CreateItemDialogProps) {
  const [name, setName] = useState("");
  const [type, setType] = useState<"file" | "folder">("file");
  const createFileApi = useFileStore((state) => state.createFileApi);
  const createFolderApi = useFileStore((state) => state.createFolderApi);
  const selectFile = useFileStore((state) => state.selectFile);

  const handleCreate = async () => {
    if (!name.trim()) return;

    try {
      if (type === "file") {
        await createFileApi(name.trim(), parentId);
      } else {
        await createFolderApi(name.trim(), parentId);
      }

      setName("");
      setType("file");
      onOpenChange(false);
    } catch (error) {
      console.error("Erro ao criar item:", error);
    }
  };

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleCreate();
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <FileText className="w-5 h-5" />
            Criar {type === "file" ? "Arquivo" : "Pasta"}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-4 py-4">
          <div className="space-y-2">
            <Label>Tipo</Label>
            <RadioGroup
              value={type}
              onValueChange={(value: "file" | "folder") => setType(value)}
            >
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="file" id="file" />
                <Label
                  htmlFor="file"
                  className="flex items-center gap-2 cursor-pointer"
                >
                  <FileText className="w-4 h-4 text-file-document" />
                  Arquivo (.excalidraw)
                </Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="folder" id="folder" />
                <Label
                  htmlFor="folder"
                  className="flex items-center gap-2 cursor-pointer"
                >
                  <Folder className="w-4 h-4 text-file-folder" />
                  Pasta
                </Label>
              </div>
            </RadioGroup>
          </div>

          <div className="space-y-2">
            <Label htmlFor="name">Nome</Label>
            <Input
              id="name"
              placeholder={type === "file" ? "Meu desenho" : "Nova pasta"}
              value={name}
              onChange={(e) => setName(e.target.value)}
              onKeyDown={handleKeyDown}
              autoFocus
            />
            {type === "file" && (
              <p className="text-xs text-muted-foreground">
                A extensão .excalidraw será adicionada automaticamente
              </p>
            )}
          </div>
        </div>

        <div className="flex justify-end gap-2">
          <Button
            variant="outline"
            onClick={() => {
              setName("");
              setType("file");
              onOpenChange(false);
            }}
          >
            Cancelar
          </Button>
          <Button
            onClick={handleCreate}
            disabled={!name.trim()}
            className="bg-primary hover:bg-primary-hover"
          >
            Criar
          </Button>
        </div>
      </DialogContent>
    </Dialog>
  );
}
