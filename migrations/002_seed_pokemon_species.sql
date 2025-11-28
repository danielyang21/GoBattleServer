-- Seed Pokemon species data
-- This includes a variety of Pokemon across different rarities

-- Mythic (extremely rare, ~0.1% drop rate)
INSERT INTO pokemon_species (id, name, rarity, base_hp, base_attack, base_defense, base_sp_attack, base_sp_defense, base_speed, sprite_url, drop_weight) VALUES
(150, 'Mewtwo', 'mythic', 106, 110, 90, 154, 90, 130, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/150.png', 0.001),
(151, 'Mew', 'mythic', 100, 100, 100, 100, 100, 100, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/151.png', 0.001),
(249, 'Lugia', 'mythic', 106, 90, 130, 90, 154, 110, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/249.png', 0.001),
(250, 'Ho-Oh', 'mythic', 106, 130, 90, 110, 154, 90, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/250.png', 0.001),
(384, 'Rayquaza', 'mythic', 105, 150, 90, 150, 90, 95, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/384.png', 0.001);

-- Legendary (~1% drop rate)
INSERT INTO pokemon_species (id, name, rarity, base_hp, base_attack, base_defense, base_sp_attack, base_sp_defense, base_speed, sprite_url, drop_weight) VALUES
(144, 'Articuno', 'legendary', 90, 85, 100, 95, 125, 85, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/144.png', 0.01),
(145, 'Zapdos', 'legendary', 90, 90, 85, 125, 90, 100, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/145.png', 0.01),
(146, 'Moltres', 'legendary', 90, 100, 90, 125, 85, 90, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/146.png', 0.01),
(243, 'Raikou', 'legendary', 90, 85, 75, 115, 100, 115, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/243.png', 0.01),
(244, 'Entei', 'legendary', 115, 115, 85, 90, 75, 100, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/244.png', 0.01),
(245, 'Suicune', 'legendary', 100, 75, 115, 90, 115, 85, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/245.png', 0.01),
(377, 'Regirock', 'legendary', 80, 100, 200, 50, 100, 50, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/377.png', 0.01),
(378, 'Regice', 'legendary', 80, 50, 100, 100, 200, 50, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/378.png', 0.01),
(379, 'Registeel', 'legendary', 80, 75, 150, 75, 150, 50, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/379.png', 0.01),
(380, 'Latias', 'legendary', 80, 80, 90, 110, 130, 110, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/380.png', 0.01);

-- Epic (~5% drop rate)
INSERT INTO pokemon_species (id, name, rarity, base_hp, base_attack, base_defense, base_sp_attack, base_sp_defense, base_speed, sprite_url, drop_weight) VALUES
(3, 'Venusaur', 'epic', 80, 82, 83, 100, 100, 80, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/3.png', 0.05),
(6, 'Charizard', 'epic', 78, 84, 78, 109, 85, 100, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/6.png', 0.05),
(9, 'Blastoise', 'epic', 79, 83, 100, 85, 105, 78, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/9.png', 0.05),
(94, 'Gengar', 'epic', 60, 65, 60, 130, 75, 110, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/94.png', 0.05),
(131, 'Lapras', 'epic', 130, 85, 80, 85, 95, 60, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/131.png', 0.05),
(143, 'Snorlax', 'epic', 160, 110, 65, 65, 110, 30, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/143.png', 0.05),
(149, 'Dragonite', 'epic', 91, 134, 95, 100, 100, 80, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/149.png', 0.05),
(248, 'Tyranitar', 'epic', 100, 134, 110, 95, 100, 61, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/248.png', 0.05),
(282, 'Gardevoir', 'epic', 68, 65, 65, 125, 115, 80, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/282.png', 0.05),
(376, 'Metagross', 'epic', 80, 135, 130, 95, 90, 70, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/376.png', 0.05);

-- Rare (~15% drop rate)
INSERT INTO pokemon_species (id, name, rarity, base_hp, base_attack, base_defense, base_sp_attack, base_sp_defense, base_speed, sprite_url, drop_weight) VALUES
(2, 'Ivysaur', 'rare', 60, 62, 63, 80, 80, 60, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/2.png', 0.15),
(5, 'Charmeleon', 'rare', 58, 64, 58, 80, 65, 80, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/5.png', 0.15),
(8, 'Wartortle', 'rare', 59, 63, 80, 65, 80, 58, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/8.png', 0.15),
(26, 'Raichu', 'rare', 60, 90, 55, 90, 80, 110, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/26.png', 0.15),
(34, 'Nidoking', 'rare', 81, 102, 77, 85, 75, 85, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/34.png', 0.15),
(59, 'Arcanine', 'rare', 90, 110, 80, 100, 80, 95, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/59.png', 0.15),
(65, 'Alakazam', 'rare', 55, 50, 45, 135, 95, 120, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/65.png', 0.15),
(68, 'Machamp', 'rare', 90, 130, 80, 65, 85, 55, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/68.png', 0.15),
(76, 'Golem', 'rare', 80, 120, 130, 55, 65, 45, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/76.png', 0.15),
(91, 'Cloyster', 'rare', 50, 95, 180, 85, 45, 70, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/91.png', 0.15),
(103, 'Exeggutor', 'rare', 95, 95, 85, 125, 75, 55, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/103.png', 0.15),
(112, 'Rhydon', 'rare', 105, 130, 120, 45, 45, 40, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/112.png', 0.15),
(130, 'Gyarados', 'rare', 95, 125, 79, 60, 100, 81, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/130.png', 0.15),
(142, 'Aerodactyl', 'rare', 80, 105, 65, 60, 75, 130, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/142.png', 0.15),
(148, 'Dragonair', 'rare', 61, 84, 65, 70, 70, 70, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/148.png', 0.15);

-- Uncommon (~30% drop rate)
INSERT INTO pokemon_species (id, name, rarity, base_hp, base_attack, base_defense, base_sp_attack, base_sp_defense, base_speed, sprite_url, drop_weight) VALUES
(1, 'Bulbasaur', 'uncommon', 45, 49, 49, 65, 65, 45, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/1.png', 0.30),
(4, 'Charmander', 'uncommon', 39, 52, 43, 60, 50, 65, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/4.png', 0.30),
(7, 'Squirtle', 'uncommon', 44, 48, 65, 50, 64, 43, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/7.png', 0.30),
(25, 'Pikachu', 'uncommon', 35, 55, 40, 50, 50, 90, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/25.png', 0.30),
(39, 'Jigglypuff', 'uncommon', 115, 45, 20, 45, 25, 20, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/39.png', 0.30),
(54, 'Psyduck', 'uncommon', 50, 52, 48, 65, 50, 55, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/54.png', 0.30),
(58, 'Growlithe', 'uncommon', 55, 70, 45, 70, 50, 60, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/58.png', 0.30),
(63, 'Abra', 'uncommon', 25, 20, 15, 105, 55, 90, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/63.png', 0.30),
(66, 'Machop', 'uncommon', 70, 80, 50, 35, 35, 35, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/66.png', 0.30),
(74, 'Geodude', 'uncommon', 40, 80, 100, 30, 30, 20, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/74.png', 0.30),
(92, 'Gastly', 'uncommon', 30, 35, 30, 100, 35, 80, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/92.png', 0.30),
(95, 'Onix', 'uncommon', 35, 45, 160, 30, 45, 70, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/95.png', 0.30),
(104, 'Cubone', 'uncommon', 50, 50, 95, 40, 50, 35, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/104.png', 0.30),
(111, 'Rhyhorn', 'uncommon', 80, 85, 95, 30, 30, 25, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/111.png', 0.30),
(133, 'Eevee', 'uncommon', 55, 55, 50, 45, 65, 55, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/133.png', 0.30),
(147, 'Dratini', 'uncommon', 41, 64, 45, 50, 50, 50, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/147.png', 0.30),
(152, 'Chikorita', 'uncommon', 45, 49, 65, 49, 65, 45, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/152.png', 0.30),
(155, 'Cyndaquil', 'uncommon', 39, 52, 43, 60, 50, 65, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/155.png', 0.30),
(158, 'Totodile', 'uncommon', 50, 65, 64, 44, 48, 43, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/158.png', 0.30),
(172, 'Pichu', 'uncommon', 20, 40, 15, 35, 35, 60, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/172.png', 0.30);

-- Common (~49% drop rate)
INSERT INTO pokemon_species (id, name, rarity, base_hp, base_attack, base_defense, base_sp_attack, base_sp_defense, base_speed, sprite_url, drop_weight) VALUES
(10, 'Caterpie', 'common', 45, 30, 35, 20, 20, 45, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/10.png', 0.49),
(13, 'Weedle', 'common', 40, 35, 30, 20, 20, 50, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/13.png', 0.49),
(16, 'Pidgey', 'common', 40, 45, 40, 35, 35, 56, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/16.png', 0.49),
(19, 'Rattata', 'common', 30, 56, 35, 25, 35, 72, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/19.png', 0.49),
(21, 'Spearow', 'common', 40, 60, 30, 31, 31, 70, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/21.png', 0.49),
(27, 'Sandshrew', 'common', 50, 75, 85, 20, 30, 40, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/27.png', 0.49),
(29, 'Nidoran-f', 'common', 55, 47, 52, 40, 40, 41, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/29.png', 0.49),
(32, 'Nidoran-m', 'common', 46, 57, 40, 40, 40, 50, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/32.png', 0.49),
(41, 'Zubat', 'common', 40, 45, 35, 30, 40, 55, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/41.png', 0.49),
(43, 'Oddish', 'common', 45, 50, 55, 75, 65, 30, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/43.png', 0.49),
(48, 'Venonat', 'common', 60, 55, 50, 40, 55, 45, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/48.png', 0.49),
(50, 'Diglett', 'common', 10, 55, 25, 35, 45, 95, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/50.png', 0.49),
(52, 'Meowth', 'common', 40, 45, 35, 40, 40, 90, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/52.png', 0.49),
(60, 'Poliwag', 'common', 40, 50, 40, 40, 40, 90, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/60.png', 0.49),
(69, 'Bellsprout', 'common', 50, 75, 35, 70, 30, 40, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/69.png', 0.49),
(72, 'Tentacool', 'common', 40, 40, 35, 50, 100, 70, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/72.png', 0.49),
(77, 'Ponyta', 'common', 50, 85, 55, 65, 65, 90, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/77.png', 0.49),
(81, 'Magnemite', 'common', 25, 35, 70, 95, 55, 45, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/81.png', 0.49),
(84, 'Doduo', 'common', 35, 85, 45, 35, 35, 75, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/84.png', 0.49),
(96, 'Drowzee', 'common', 60, 48, 45, 43, 90, 42, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/96.png', 0.49),
(98, 'Krabby', 'common', 30, 105, 90, 25, 25, 50, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/98.png', 0.49),
(100, 'Voltorb', 'common', 40, 30, 50, 55, 55, 100, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/100.png', 0.49),
(109, 'Koffing', 'common', 40, 65, 95, 60, 45, 35, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/109.png', 0.49),
(118, 'Goldeen', 'common', 45, 67, 60, 35, 50, 63, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/118.png', 0.49),
(120, 'Staryu', 'common', 30, 45, 55, 70, 55, 85, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/120.png', 0.49),
(129, 'Magikarp', 'common', 20, 10, 55, 15, 20, 80, 'https://raw.githubusercontent.com/PokeAPI/sprites/master/sprites/pokemon/129.png', 0.49);
