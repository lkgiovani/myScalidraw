import {
  ContextMenu,
  ContextMenuContent,
  ContextMenuTrigger,
  ContextMenuItem,
  ContextMenuSeparator,
} from "@/components/ui/context-menu";
import { useFileStore, FileNode } from "@/stores/fileStore";
import { useDeleteFile } from "@/hooks/useFileApi";
import { useState } from "react";
import { CreateItemDialog } from "./CreateItemDialog";
import { RenameDialog } from "./RenameDialog";
import { FileText, Folder } from "lucide-react";

interface FileContextMenuProps {
  file?: FileNode;
  onCreateNew?: () => void;
  children: React.ReactNode;
}

export function FileContextMenu({
  file,
  onCreateNew,
  children,
}: FileContextMenuProps) {
  const deleteItem = useFileStore((state) => state.deleteItem);
  const deleteFileApi = useFileStore((state) => state.deleteFileApi);
  const deleteFileMutation = useDeleteFile();
  const [showCreateDialog, setShowCreateDialog] = useState(false);
  const [showRenameDialog, setShowRenameDialog] = useState(false);

  const handleDelete = async () => {
    if (file && confirm(`Tem certeza que deseja deletar "${file.name}"?`)) {
      try {
        await deleteFileMutation.mutateAsync(file.id);

        deleteItem(file.id);
      } catch (error) {
        console.error("Erro ao deletar arquivo:", error);

        deleteItem(file.id);
      }
    }
  };

  const handleCreateNew = () => {
    if (onCreateNew) {
      onCreateNew();
    } else {
      setShowCreateDialog(true);
    }
  };

  return (
    <>
      <ContextMenu>
        <ContextMenuTrigger asChild>{children}</ContextMenuTrigger>
        <ContextMenuContent className="w-48">
          <ContextMenuItem onClick={handleCreateNew} className="gap-2">
            <FileText className="w-4 h-4" />
            Novo Item
          </ContextMenuItem>

          {file?.type === "folder" && (
            <ContextMenuItem
              onClick={() => setShowCreateDialog(true)}
              className="gap-2"
            >
              <Folder className="w-4 h-4" />
              Novo em "{file.name}"
            </ContextMenuItem>
          )}

          {file && (
            <>
              <ContextMenuSeparator />
              <ContextMenuItem
                onClick={() => setShowRenameDialog(true)}
                className="gap-2"
              >
                <FileText className="w-4 h-4" />
                Renomear
              </ContextMenuItem>
              <ContextMenuItem
                onClick={handleDelete}
                className="gap-2 text-destructive focus:text-destructive"
              >
                <FileText className="w-4 h-4" />
                Deletar
              </ContextMenuItem>
            </>
          )}
        </ContextMenuContent>
      </ContextMenu>

      <CreateItemDialog
        open={showCreateDialog}
        onOpenChange={setShowCreateDialog}
        parentId={file?.type === "folder" ? file.id : undefined}
      />

      {file && (
        <RenameDialog
          open={showRenameDialog}
          onOpenChange={setShowRenameDialog}
          file={file}
        />
      )}
    </>
  );
}
