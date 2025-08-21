import React from "react";
import { Link } from "react-router-dom";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

const Home = () => {
  return (
    <div className="min-h-screen bg-gradient-to-br from-background via-background to-card">
      <div className="container mx-auto px-4 py-16">
        {/* Header */}
        <div className="text-center mb-16">
          <div className="w-24 h-24 mx-auto mb-8 bg-primary/20 rounded-full flex items-center justify-center">
            <svg
              className="w-12 h-12 text-primary"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"
              />
            </svg>
          </div>
          <h1 className="text-5xl font-bold text-foreground mb-4">
            myScalidraw
          </h1>
          <p className="text-xl text-muted-foreground max-w-2xl mx-auto">
            Uma ferramenta de desenho colaborativa e intuitiva baseada no
            Excalidraw. Crie diagramas, wireframes e ilustrações de forma
            simples e eficiente.
          </p>
        </div>

        {/* Features Grid */}
        <div className="grid md:grid-cols-3 gap-8 mb-16">
          <Card className="border-border/50 hover:border-primary/50 transition-colors">
            <CardHeader>
              <div className="w-12 h-12 bg-primary/20 rounded-lg flex items-center justify-center mb-4">
                <svg
                  className="w-6 h-6 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M7 21a4 4 0 01-4-4V5a2 2 0 012-2h4a2 2 0 012 2v12a4 4 0 01-4 4zM21 5a2 2 0 00-2-2h-4a2 2 0 00-2 2v6a4 4 0 004 4h4V5z"
                  />
                </svg>
              </div>
              <CardTitle>Interface Intuitiva</CardTitle>
              <CardDescription>
                Interface limpa e moderna com tema escuro profissional,
                otimizada para produtividade.
              </CardDescription>
            </CardHeader>
          </Card>

          <Card className="border-border/50 hover:border-primary/50 transition-colors">
            <CardHeader>
              <div className="w-12 h-12 bg-primary/20 rounded-lg flex items-center justify-center mb-4">
                <svg
                  className="w-6 h-6 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"
                  />
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M8 1v6m8-6v6"
                  />
                </svg>
              </div>
              <CardTitle>Gerenciamento de Arquivos</CardTitle>
              <CardDescription>
                Organize seus desenhos em pastas, renomeie e gerencie seus
                projetos de forma eficiente.
              </CardDescription>
            </CardHeader>
          </Card>

          <Card className="border-border/50 hover:border-primary/50 transition-colors">
            <CardHeader>
              <div className="w-12 h-12 bg-primary/20 rounded-lg flex items-center justify-center mb-4">
                <svg
                  className="w-6 h-6 text-primary"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M13 10V3L4 14h7v7l9-11h-7z"
                  />
                </svg>
              </div>
              <CardTitle>Performance Otimizada</CardTitle>
              <CardDescription>
                Baseado no Excalidraw com otimizações para uma experiência
                fluida e responsiva.
              </CardDescription>
            </CardHeader>
          </Card>
        </div>

        {/* CTA Section */}
        <div className="text-center">
          <Card className="max-w-2xl mx-auto border-primary/20 bg-card/50">
            <CardContent className="pt-8">
              <h2 className="text-3xl font-bold text-foreground mb-4">
                Pronto para começar?
              </h2>
              <p className="text-muted-foreground mb-8">
                Acesse o editor e comece a criar seus desenhos e diagramas agora
                mesmo.
              </p>
              <div className="flex gap-4 justify-center">
                <Link to="/editor">
                  <Button size="lg" className="px-8">
                    Abrir Editor
                  </Button>
                </Link>
                <Button variant="outline" size="lg" className="px-8">
                  Ver Exemplos
                </Button>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Footer */}
        <div className="text-center mt-16 pt-8 border-t border-border/50">
          <p className="text-muted-foreground">
            Desenvolvido com React, TypeScript e Excalidraw
          </p>
        </div>
      </div>
    </div>
  );
};

export default Home;
