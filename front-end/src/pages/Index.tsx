import React from 'react';
import { FileExplorer } from '@/components/FileExplorer';
import { ExcalidrawCanvas } from '@/components/ExcalidrawCanvas';
import { AppHeader } from '@/components/AppHeader';
import { useFileSystem } from '@/hooks/useFileSystem';
import { toast } from 'sonner';

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
  } = useFileSystem();

  const activeFile = getActiveFile();

  const handleFileSelect = (id: string) => {
    setActiveFileId(id);
  };

  const handleCreateFile = (name: string, parentId?: string) => {
    const fileId = createFile(name, parentId);
    toast.success(`Arquivo "${name}" criado com sucesso`);
    return fileId;
  };

  const handleCreateFolder = (name: string, parentId?: string) => {
    const folderId = createFolder(name, parentId);
    toast.success(`Pasta "${name}" criada com sucesso`);
    return folderId;
  };

  const handleRename = (id: string, newName: string) => {
    renameItem(id, newName);
    toast.success('Item renomeado com sucesso');
  };

  const handleDelete = (id: string) => {
    deleteItem(id);
    toast.success('Item excluído com sucesso');
  };

  const handleSave = (id: string, data: any) => {
    updateFile(id, data);
  };

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
          <ExcalidrawCanvas
            file={activeFile}
            onSave={handleSave}
          />
        </div>
      </div>
    </div>
  );
};

export default Index;
