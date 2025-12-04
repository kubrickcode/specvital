import { UrlInputForm } from "@/components/url-input-form";

export const dynamic = "force-static";

export const HomePage = () => {
  return (
    <main className="flex min-h-screen flex-col items-center justify-center p-8">
      <div className="w-full max-w-2xl space-y-8 text-center">
        <div className="space-y-4">
          <h1 className="text-4xl font-bold tracking-tight sm:text-5xl">SpecVital</h1>
          <p className="text-lg text-muted-foreground">
            Analyze test specifications from public GitHub repositories
          </p>
        </div>

        <UrlInputForm />

        <p className="text-sm text-muted-foreground">
          Enter a public GitHub repository URL to analyze its test files
        </p>
      </div>
    </main>
  );
};

export default HomePage;
