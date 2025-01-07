import "jsr:@supabase/functions-js/edge-runtime.d.ts";
import { getApiClient } from "https://deno.land/x/pokemontcg_deno@v0.1.0/mod.ts";

Deno.serve(async (req) => {
  const TCG_API_KEY = Deno.env.get("TCG_API_KEY");

  const client = getApiClient(TCG_API_KEY);

  const cards = await client.cards.all({ pageSize: 100 });

  return new Response(JSON.stringify({ cards: cards.length }), {
    headers: { "Content-Type": "application/json" },
  });
});
