import React, { useState } from "react";
import {
  ChevronDown,
  ChevronRight,
  FolderPlus,
  FilePlus,
  MoreHorizontal,
  File,
  Folder,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  FileEntity,
  FileSystemItemEntity,
  FolderEntity,
} from "@/domain/entities/FileEntity";
import { cn } from "@/lib/utils";

interface FileExplorerProps {
  files: FileSystemItemEntity[];
  activeFileId: string;
  onFileSelect: (id: string) => void;
  onCreateFile: (name: string, parentId?: string) => string | Promise<string>;
  onCreateFolder: (name: string, parentId?: string) => string | Promise<string>;
  onRename: (id: string, newName: string) => void;
  onDelete: (id: string) => void;
  onToggleFolder: (id: string) => void;
}

interface FileItemProps {
  item: FileSystemItemEntity;
  level: number;
  isActive: boolean;
  activeFileId: string;
  onSelect: (id: string) => void;
  onRename: (id: string, newName: string) => void;
  onDelete: (id: string) => void;
  onToggleFolder: (id: string) => void;
}

const FileItem: React.FC<FileItemProps> = ({
  item,
  level,
  isActive,
  activeFileId,
  onSelect,
  onRename,
  onDelete,
  onToggleFolder,
}) => {
  const [isRenaming, setIsRenaming] = useState(false);
  const [renameName, setRenameName] = useState(item.name);

  const handleRename = () => {
    if (renameName.trim() && renameName !== item.name) {
      onRename(item.id, renameName.trim());
    }
    setIsRenaming(false);
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      handleRename();
    } else if (e.key === "Escape") {
      setIsRenaming(false);
      setRenameName(item.name);
    }
  };

  const isFolder = "isFolder" in item;
  const paddingLeft = level * 16 + 8;

  return (
    <div className="select-none">
      <div
        className={cn(
          "group flex items-center h-8 px-2 hover:bg-file-hover cursor-pointer transition-smooth relative",
          isActive &&
            !isFolder &&
            "bg-primary/20 border-r-2 border-primary text-primary-foreground"
        )}
        style={{ paddingLeft }}
        onClick={() => {
          if (isFolder) {
            onToggleFolder(item.id);
          } else {
            onSelect(item.id);
          }
        }}
      >
        {isFolder ? (
          <>
            {(item as FolderEntity).isExpanded ? (
              <ChevronDown className="h-4 w-4 mr-1 text-muted-foreground" />
            ) : (
              <ChevronRight className="h-4 w-4 mr-1 text-muted-foreground" />
            )}
            <Folder className="h-4 w-4 mr-2 text-folder-icon" />
          </>
        ) : (
          <>
            <div className="w-5 mr-1" />
            <File className="h-4 w-4 mr-2 text-file-icon" />
          </>
        )}

        {isRenaming ? (
          <Input
            value={renameName}
            onChange={(e) => setRenameName(e.target.value)}
            onBlur={handleRename}
            onKeyDown={handleKeyPress}
            className="h-6 text-xs border-primary focus-ring"
            autoFocus
            onClick={(e) => e.stopPropagation()}
          />
        ) : (
          <span className="text-sm truncate flex-1">{item.name}</span>
        )}

        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button
              variant="ghost"
              size="sm"
              className={cn(
                "h-6 w-6 p-0 opacity-0 group-hover:opacity-100 transition-smooth hover:bg-accent",
                isActive && "opacity-100"
              )}
              onClick={(e) => e.stopPropagation()}
            >
              <MoreHorizontal className="h-3 w-3" />
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-40">
            <DropdownMenuItem onClick={() => setIsRenaming(true)}>
              Renomear
            </DropdownMenuItem>
            <DropdownMenuItem
              onClick={() => onDelete(item.id)}
              className="text-destructive focus:text-destructive"
            >
              Excluir
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>

      {isFolder && (item as FolderEntity).isExpanded && (
        <div className="animate-fade-in">
          {(item as FolderEntity).children.map((child) => (
            <FileItem
              key={child.id}
              item={child}
              level={level + 1}
              isActive={child.id === activeFileId}
              activeFileId={activeFileId}
              onSelect={onSelect}
              onRename={onRename}
              onDelete={onDelete}
              onToggleFolder={onToggleFolder}
            />
          ))}
        </div>
      )}
    </div>
  );
};

export const FileExplorer: React.FC<FileExplorerProps> = ({
  files,
  activeFileId,
  onFileSelect,
  onCreateFile,
  onCreateFolder,
  onRename,
  onDelete,
  onToggleFolder,
}) => {
  const [isCreatingFile, setIsCreatingFile] = useState(false);
  const [isCreatingFolder, setIsCreatingFolder] = useState(false);
  const [newItemName, setNewItemName] = useState("");

  const renderFileTree = (
    items: FileSystemItemEntity[],
    level: number = 0
  ): React.ReactNode => {
    return items.map((item) => (
      <FileItem
        key={item.id}
        item={item}
        level={level}
        isActive={item.id === activeFileId}
        activeFileId={activeFileId}
        onSelect={onFileSelect}
        onRename={onRename}
        onDelete={onDelete}
        onToggleFolder={onToggleFolder}
      />
    ));
  };

  const handleCreateFile = async () => {
    if (newItemName.trim()) {
      await onCreateFile(newItemName.trim());
      setNewItemName("");
      setIsCreatingFile(false);
    }
  };

  const handleCreateFolder = async () => {
    if (newItemName.trim()) {
      await onCreateFolder(newItemName.trim());
      setNewItemName("");
      setIsCreatingFolder(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === "Enter") {
      if (isCreatingFile) {
        handleCreateFile();
      } else if (isCreatingFolder) {
        handleCreateFolder();
      }
    } else if (e.key === "Escape") {
      setIsCreatingFile(false);
      setIsCreatingFolder(false);
      setNewItemName("");
    }
  };

  return (
    <div className="h-full flex flex-col bg-sidebar-background border-r border-sidebar-border">
      <div className="p-3 border-b border-sidebar-border">
        <div className="flex items-center justify-between mb-2">
          <h2 className="text-sm font-semibold text-sidebar-foreground">
            Explorador
          </h2>
          <div className="flex gap-1">
            <Button
              variant="ghost"
              size="sm"
              className="h-6 w-6 p-0 hover:bg-sidebar-accent"
              onClick={() => setIsCreatingFile(true)}
              title="Novo arquivo"
            >
              <FilePlus className="h-3 w-3" />
            </Button>
            <Button
              variant="ghost"
              size="sm"
              className="h-6 w-6 p-0 hover:bg-sidebar-accent"
              onClick={() => setIsCreatingFolder(true)}
              title="Nova pasta"
            >
              <FolderPlus className="h-3 w-3" />
            </Button>
          </div>
        </div>

        {(isCreatingFile || isCreatingFolder) && (
          <div className="mb-2">
            <Input
              value={newItemName}
              onChange={(e) => setNewItemName(e.target.value)}
              onBlur={isCreatingFile ? handleCreateFile : handleCreateFolder}
              onKeyDown={handleKeyPress}
              placeholder={isCreatingFile ? "Nome do arquivo" : "Nome da pasta"}
              className="h-7 text-xs border-primary focus-ring"
              autoFocus
            />
          </div>
        )}
      </div>

      <div className="flex-1 overflow-auto">
        <div className="animate-fade-in">{renderFileTree(files)}</div>

        {files.length === 0 && (
          <div className="p-4 text-center text-muted-foreground text-sm">
            <p>Nenhum arquivo ainda.</p>
            <p>Clique nos ícones acima para criar.</p>
          </div>
        )}
      </div>
    </div>
  );
};
