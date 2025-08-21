import React from 'react';
import { ExcalidrawFile } from '@/types/file';

interface AppHeaderProps {
  activeFile: ExcalidrawFile | undefined;
}

export const AppHeader: React.FC<AppHeaderProps> = ({ activeFile }) => {
  return (
    <header className="h-12 bg-card border-b border-border flex items-center px-4">
      <div className="flex items-center gap-3">
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 rounded-full bg-destructive"></div>
          <div className="w-3 h-3 rounded-full bg-yellow-500"></div>
          <div className="w-3 h-3 rounded-full bg-green-500"></div>
        </div>
        
        <div className="w-px h-6 bg-border"></div>
        
        <h1 className="text-sm font-medium text-foreground">
          React Excalidraw Editor
        </h1>
        
        {activeFile && (
          <>
            <div className="w-px h-6 bg-border"></div>
            <span className="text-sm text-muted-foreground">
              {activeFile.name}
            </span>
          </>
        )}
      </div>
    </header>
  );
};