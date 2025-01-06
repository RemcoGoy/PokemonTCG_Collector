import "jsr:@supabase/functions-js/edge-runtime.d.ts";
import { getCards } from "https://x.nest.land/pokedeno@0.2.1/mod.ts";

Deno.serve(async (req) => {
  const TCG_API_KEY = Deno.env.get("TCG_API_KEY");

  const cards = await getCards({});
  const ids = cards.map((card) => card.id);

  return new Response(JSON.stringify({ ids }), {
    headers: { "Content-Type": "application/json" },
  });
});
