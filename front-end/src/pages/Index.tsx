import React, { useEffect, useState } from "react";
import { FileExplorer } from "@/components/FileExplorer";
import { ExcalidrawCanvas } from "@/components/ExcalidrawCanvas";
import { AppHeader } from "@/components/AppHeader";
import { useFileSystemAdapter } from "@/hooks/useFileSystemAdapter";
import { toast } from "sonner";
import { FileEntity } from "@/domain/entities/FileEntity";

const Index = () => {
  const {
    files,
    activeFileId,
    setActiveFileId,
    createFile,
    createFolder,
    updateFile,
    renameItem,
    deleteItem,
    toggleFolder,
    getActiveFile,
    isLoading,
  } = useFileSystemAdapter();

  const [activeFile, setActiveFile] = useState<FileEntity | undefined>(
    undefined
  );

  // Get active file when activeFileId changes
  useEffect(() => {
    const fetchActiveFile = async () => {
      if (activeFileId) {
        const file = await getActiveFile();
        setActiveFile(file);
      } else {
        setActiveFile(undefined);
      }
    };

    fetchActiveFile();
  }, [activeFileId, getActiveFile]);

  const handleFileSelect = (id: string) => {
    setActiveFileId(id);
  };

  const handleCreateFile = async (name: string, parentId?: string) => {
    try {
      const fileId = await createFile(name, parentId);
      toast.success(`Arquivo "${name}" criado com sucesso`);
      return fileId;
    } catch (error) {
      toast.error("Erro ao criar arquivo");
      return "";
    }
  };

  const handleCreateFolder = async (name: string, parentId?: string) => {
    try {
      const folderId = await createFolder(name, parentId);
      toast.success(`Pasta "${name}" criada com sucesso`);
      return folderId;
    } catch (error) {
      toast.error("Erro ao criar pasta");
      return "";
    }
  };

  const handleRename = async (id: string, newName: string) => {
    try {
      await renameItem(id, newName);
      toast.success("Item renomeado com sucesso");
    } catch (error) {
      toast.error("Erro ao renomear item");
    }
  };

  const handleDelete = async (id: string) => {
    try {
      await deleteItem(id);
      toast.success("Item excluído com sucesso");
    } catch (error) {
      toast.error("Erro ao excluir item");
    }
  };

  const handleSave = async (id: string, data: any) => {
    try {
      await updateFile(id, data);
    } catch (error) {
      toast.error("Erro ao salvar arquivo");
    }
  };

  if (isLoading) {
    return (
      <div className="h-screen flex items-center justify-center bg-background">
        <div className="text-center">
          <div className="w-16 h-16 border-4 border-primary border-t-transparent rounded-full animate-spin mx-auto mb-4"></div>
          <p className="text-muted-foreground">Carregando arquivos...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="h-screen flex flex-col bg-background text-foreground overflow-hidden">
      <AppHeader activeFile={activeFile} />

      <div className="flex-1 flex overflow-hidden">
        <div className="w-80 flex-shrink-0">
          <FileExplorer
            files={files}
            activeFileId={activeFileId}
            onFileSelect={handleFileSelect}
            onCreateFile={handleCreateFile}
            onCreateFolder={handleCreateFolder}
            onRename={handleRename}
            onDelete={handleDelete}
            onToggleFolder={toggleFolder}
          />
        </div>

        <div className="flex-1 overflow-hidden">
          <ExcalidrawCanvas file={activeFile} onSave={handleSave} />
        </div>
      </div>
    </div>
  );
};

export default Index;
