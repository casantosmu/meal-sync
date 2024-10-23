CREATE TABLE recipes (
    recipe_id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    img_url TEXT NOT NULL,
    description TEXT,
    ingredients TEXT,
    directions TEXT
);

CREATE INDEX idx_recipes_title ON recipes(title, ingredients);

INSERT INTO recipes (title, img_url, description, ingredients, directions)
VALUES
  ('Spicy Chorizo and Egg Burgers with White Pepper', 
   'https://placehold.co/500x500', 
   'A deliciously spicy twist on a classic burger, combining crispy chorizo with eggs fried to perfection. The addition of white pepper adds a subtle kick that makes this burger stand out. Perfect for brunch or a satisfying dinner.',
   'Chorizo - 200g
Eggs - 4
White pepper, a pinch
Buns
Cheese slices
Lettuce
Tomato',
   '1. Heat a skillet over medium-high heat and cook the chorizo until crispy, about 5-7 minutes. Set aside.

2. In the same pan, fry the eggs to your liking (sunny side up, over-easy, or well-done). Sprinkle with a pinch of white pepper for an extra layer of flavor.

3. Lightly toast the buns in the same skillet.

4. Assemble the burgers by layering the crispy chorizo, fried eggs, cheese slices, lettuce, and tomato between the toasted buns.

5. Serve immediately while warm and enjoy with your favorite side.'
  ),
  
  ('Quinoa Salad with Mango and Avocado', 
   'https://placehold.co/500x500', 
   'This refreshing quinoa salad combines the tropical sweetness of mango with the creamy texture of avocado. Perfect as a light lunch or a healthy side dish, it is bursting with vibrant flavors and colors.',
   '1 cup Quinoa
Mango (diced)
Avocado - 1
Red onion, chopped
2 tbsp Lime juice
Olive oil
Salt and pepper',
   '1. Rinse the quinoa thoroughly under cold water, then cook according to the package instructions, typically by simmering for 15 minutes until the water is absorbed and the quinoa is fluffy. Let it cool.

2. In a large bowl, combine the cooled quinoa, diced mango, chopped avocado, and finely chopped red onion.

3. Drizzle the lime juice and olive oil over the mixture and gently toss to combine.

4. Season with salt and pepper to taste.

5. Serve chilled or at room temperature as a side or main dish.'
  ),
  
  ('Roasted Tomato Soup with Basil', 
   'https://placehold.co/500x500', 
   'A comforting and flavorful roasted tomato soup, enhanced with the sweetness of roasted garlic and the freshness of basil. This rich soup is perfect for cold days and can be served with crusty bread for dipping.',
   '6 Tomatoes
Garlic
Olive oil
Fresh basil
Vegetable broth
Salt and pepper
Cream (optional)',
   '1. Preheat the oven to 200°C. Place the tomatoes and garlic cloves on a baking sheet, drizzle with olive oil, and season with salt and pepper. Roast for 25-30 minutes until the tomatoes are blistered and garlic is soft.

2. Once roasted, transfer the tomatoes and garlic to a blender. Add the vegetable broth and blend until smooth.

3. Pour the soup into a pot and bring it to a simmer over medium heat. Stir in chopped fresh basil and adjust seasoning if necessary.

4. For a creamier version, add a splash of cream and stir well.

5. Serve warm, garnished with additional basil leaves if desired.'
  ),
  
  ('Baja Style Fish Tacos', 
   'https://placehold.co/500x500', 
   'Crispy, flavorful fish wrapped in warm tortillas with a fresh, tangy cabbage slaw. These Baja style tacos are light, refreshing, and perfect for a summer meal or a casual get-together.',
   '500g White fish
Tortillas
Cabbage, shredded
Lime juice
Cilantro leaves
Mayonnaise
Salsa
Avocado',
   '1. Season the white fish fillets with salt, pepper, and a splash of lime juice. Preheat a grill or a pan over medium heat.

2. Cook the fish for about 3-4 minutes per side until the flesh is opaque and flakes easily with a fork.

3. In a separate bowl, mix the shredded cabbage with a drizzle of lime juice, a tablespoon of mayonnaise, and some chopped cilantro leaves for a quick slaw.

4. Warm the tortillas on the grill or in a pan, then assemble the tacos by placing the fish fillets inside each tortilla. Top with the cabbage slaw, salsa, and sliced avocado.

5. Serve immediately with extra lime wedges on the side.'
  ),
  
  ('Classic Margherita Pizza with Fresh Mozzarella', 
   'https://placehold.co/500x500', 
   'A timeless pizza classic that highlights simple ingredients like fresh mozzarella, vibrant basil, and rich tomato sauce. This pizza is a perfect balance of flavors and textures, ideal for any pizza night.',
   'Pizza dough
Tomato sauce
Fresh mozzarella - 200g
Basil leaves
Olive oil
Salt',
   '1. Preheat your oven to 220°C (425°F) and prepare a baking sheet or pizza stone.

2. Roll out the pizza dough to your desired thickness on a lightly floured surface.

3. Spread an even layer of tomato sauce over the dough, leaving a small border around the edges for the crust.

4. Tear the fresh mozzarella into pieces and distribute them evenly over the sauce.

5. Add fresh basil leaves on top and drizzle with a little olive oil. Sprinkle with salt for added flavor.

6. Bake the pizza for 10-12 minutes, or until the crust is golden and the cheese is bubbly and slightly browned.

7. Remove from the oven and let it cool for a minute before slicing. Serve hot.'
  ),
  
  ('Traditional Seafood Paella', 
   'https://placehold.co/500x500', 
   'A vibrant and aromatic Spanish dish filled with a medley of fresh seafood and saffron-infused rice. This classic paella is packed with flavors from the chorizo, bell peppers, and a mix of shellfish, making it a showstopper at any dinner table.',
   '2 cups Rice
Shrimp - 300g
Clams (cleaned)
Mussels (cleaned)
Chorizo - 150g
Saffron, a pinch
Paprika
Olive oil
Garlic
Onion, diced
Bell peppers - 2
Peas',
   '1. Heat olive oil in a large, deep pan or paella pan over medium heat. Add diced chorizo and cook for about 5 minutes until browned and crispy.

2. Add the chopped garlic, onion, and bell peppers to the pan and sauté until the vegetables are soft and fragrant, about 5 minutes.

3. Stir in the rice, paprika, and a pinch of saffron, coating the rice in the oil and allowing it to toast slightly.

4. Add enough broth (or water) to cover the rice, and bring it to a boil. Lower the heat and simmer for about 10 minutes, stirring occasionally.

5. Add the shrimp, clams, mussels, and peas to the pan. Cover and cook for another 10-12 minutes until the rice is fully cooked and the seafood is done.

6. Let the paella rest for a few minutes before serving. Garnish with fresh lemon wedges.'
  ),
  
  ('Shrimp Ceviche with Lime and Cilantro', 
   'https://placehold.co/500x500', 
   'A zesty and refreshing ceviche that lets the natural sweetness of shrimp shine, balanced with tangy lime and fresh cilantro. This dish is perfect as an appetizer or a light summer meal, served chilled.',
   'Shrimp (peeled and deveined) - 200g
1/2 cup Lime juice
Cilantro (chopped)
Red onion, finely chopped
Jalapeño, minced
Tomatoes (diced)
Avocado
Salt and pepper',
   '1. In a medium bowl, toss the shrimp with lime juice, ensuring they are fully submerged. The acid in the lime juice will "cook" the shrimp.

2. Let the shrimp sit in the lime juice for 15-20 minutes, stirring occasionally until the shrimp turn pink and opaque.

3. Once the shrimp are cooked, drain excess lime juice and add chopped cilantro, finely chopped red onion, minced jalapeño, diced tomatoes, and cubed avocado.

4. Season with salt and pepper to taste. Stir gently to combine all ingredients.

5. Serve chilled, accompanied by tortilla chips or lettuce leaves for wrapping.'
  ),
  
  ('Coconut Curry Chicken', 
   'https://placehold.co/500x500', 
   'A creamy and aromatic coconut curry chicken dish, rich in spices and balanced with the sweetness of coconut milk. This easy-to-make curry pairs perfectly with rice for a satisfying and comforting meal.',
   'Chicken - 500g
1 can Coconut milk
Curry powder - 2 tbsp
Garlic - 3 cloves
Ginger - 1 inch
Onion
Cilantro
Lime juice
Serve with Rice',
   '1. Heat a large pan over medium heat and add a tablespoon of olive oil. Sauté the finely chopped onion, minced garlic, and grated ginger until softened and fragrant, about 5 minutes.

2. Add the chicken pieces to the pan and sprinkle with curry powder. Cook the chicken until browned on all sides, about 6-8 minutes.

3. Pour in the coconut milk, reduce the heat to low, and let the curry simmer for 20 minutes, allowing the flavors to meld and the chicken to cook through.

4. Just before serving, stir in fresh cilantro and a splash of lime juice to brighten the dish.

5. Serve the coconut curry chicken over steamed rice for a complete meal.'
  ),
  
  ('Spinach and Ricotta Ravioli with Tomato Sauce', 
   'https://placehold.co/500x500', 
   'Delicate spinach and ricotta-filled ravioli served with a simple yet flavorful tomato sauce. This comforting Italian dish is perfect for a cozy dinner, topped with fresh basil and parmesan cheese.',
   'Spinach and Ricotta Ravioli
Tomatoes, crushed - 400g
Garlic cloves, minced
Olive oil
Fresh basil
Parmesan cheese
Salt and pepper',
   '1. Bring a large pot of salted water to a boil. Add the spinach and ricotta ravioli and cook according to package instructions, usually 4-5 minutes, until they float to the top.

2. While the ravioli are cooking, heat olive oil in a saucepan over medium heat. Add minced garlic and sauté for 1-2 minutes until fragrant.

3. Add the crushed tomatoes to the pan and bring to a simmer. Season with salt, pepper, and fresh basil.

4. Once the ravioli are cooked, drain them and toss gently in the tomato sauce.

5. Serve with freshly grated parmesan cheese and additional basil leaves.'
  ),
  
  ('Chocolate Walnut Brownies', 
   'https://placehold.co/500x500', 
   'Rich and fudgy chocolate brownies studded with crunchy walnuts. These decadent brownies are a perfect dessert for chocolate lovers, with a crispy top and gooey center.',
   '200g Chocolate
Butter - 100g
1 cup Sugar
Eggs - 3
Flour - 1/2 cup
Walnuts, chopped
Vanilla extract
Salt',
   '1. Preheat your oven to 180°C (350°F) and grease a baking pan.

2. In a heatproof bowl, melt the chocolate and butter together over a pot of simmering water (double boiler method), stirring until smooth.

3. Remove the bowl from the heat and stir in the sugar until fully incorporated.

4. Add the eggs, one at a time, beating well after each addition. Stir in the vanilla extract.

5. Sift in the flour and a pinch of salt, then fold gently until just combined. Stir in the chopped walnuts.

6. Pour the batter into the prepared baking pan and smooth the top. Bake for 25-30 minutes, or until a toothpick inserted into the center comes out with a few moist crumbs.

7. Allow the brownies to cool before cutting into squares and serving.'
  );
