-- Migration: Seed essential moves for battle system
-- This includes ~50 core moves covering all types and strategies

-- =====================================================
-- NORMAL TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Tackle', 'normal', 'physical', 40, 100, 35, 0, 'A physical attack in which the user charges and slams into the target with its whole body.'),
('Scratch', 'normal', 'physical', 40, 100, 35, 0, 'Hard, pointed, sharp claws rake the target to inflict damage.'),
('Quick Attack', 'normal', 'physical', 40, 100, 30, 1, 'The user lunges at the target at a speed that makes it almost invisible. This move always goes first.'),
('Body Slam', 'normal', 'physical', 85, 100, 15, 0, 'The user drops onto the target with its full body weight. This may also leave the target with paralysis.');

-- Update Body Slam with secondary effect
UPDATE moves SET secondary_effect = '{"chance": 30, "status_inflict": {"status": "paralysis", "chance": 30}}'::jsonb WHERE name = 'Body Slam';

-- =====================================================
-- FIRE TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Ember', 'fire', 'special', 40, 100, 25, 0, 'The target is attacked with small flames. This may also leave the target with a burn.'),
('Flamethrower', 'fire', 'special', 90, 100, 15, 0, 'The target is scorched with an intense blast of fire. This may also leave the target with a burn.'),
('Fire Blast', 'fire', 'special', 110, 85, 5, 0, 'The target is attacked with an intense blast of all-consuming fire. This may also leave the target with a burn.'),
('Will-O-Wisp', 'fire', 'status', NULL, 85, 15, 0, 'The user shoots a sinister, bluish-white flame at the target to inflict a burn.');

UPDATE moves SET secondary_effect = '{"chance": 10, "status_inflict": {"status": "burn", "chance": 10}}'::jsonb WHERE name = 'Ember';
UPDATE moves SET secondary_effect = '{"chance": 10, "status_inflict": {"status": "burn", "chance": 10}}'::jsonb WHERE name = 'Flamethrower';
UPDATE moves SET secondary_effect = '{"chance": 10, "status_inflict": {"status": "burn", "chance": 10}}'::jsonb WHERE name = 'Fire Blast';
UPDATE moves SET status_inflict = '{"status": "burn", "chance": 100}'::jsonb WHERE name = 'Will-O-Wisp';

-- =====================================================
-- WATER TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Water Gun', 'water', 'special', 40, 100, 25, 0, 'The target is blasted with a forceful shot of water.'),
('Bubble Beam', 'water', 'special', 65, 100, 20, 0, 'A spray of bubbles is forcefully ejected at the target. This may also lower the target''s Speed stat.'),
('Surf', 'water', 'special', 90, 100, 15, 0, 'The user attacks everything around it by swamping its surroundings with a giant wave.'),
('Hydro Pump', 'water', 'special', 110, 80, 5, 0, 'The target is blasted by a huge volume of water launched under great pressure.'),
('Aqua Jet', 'water', 'physical', 40, 100, 20, 1, 'The user lunges at the target at a speed that makes it almost invisible. This move always goes first.');

UPDATE moves SET secondary_effect = '{"chance": 10, "stat_changes": [{"stat": "speed", "stages": -1, "target": "opponent"}]}'::jsonb WHERE name = 'Bubble Beam';

-- =====================================================
-- ELECTRIC TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Thunder Shock', 'electric', 'special', 40, 100, 30, 0, 'A jolt of electricity crashes down on the target to inflict damage. This may also leave the target with paralysis.'),
('Thunderbolt', 'electric', 'special', 90, 100, 15, 0, 'A strong electric blast crashes down on the target. This may also leave the target with paralysis.'),
('Thunder', 'electric', 'special', 110, 70, 10, 0, 'A wicked thunderbolt is dropped on the target to inflict damage. This may also leave the target with paralysis.'),
('Thunder Wave', 'electric', 'status', NULL, 90, 20, 0, 'The user launches a weak jolt of electricity that paralyzes the target.');

UPDATE moves SET secondary_effect = '{"chance": 10, "status_inflict": {"status": "paralysis", "chance": 10}}'::jsonb WHERE name = 'Thunder Shock';
UPDATE moves SET secondary_effect = '{"chance": 10, "status_inflict": {"status": "paralysis", "chance": 10}}'::jsonb WHERE name = 'Thunderbolt';
UPDATE moves SET secondary_effect = '{"chance": 30, "status_inflict": {"status": "paralysis", "chance": 30}}'::jsonb WHERE name = 'Thunder';
UPDATE moves SET status_inflict = '{"status": "paralysis", "chance": 100}'::jsonb WHERE name = 'Thunder Wave';

-- =====================================================
-- GRASS TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Vine Whip', 'grass', 'physical', 45, 100, 25, 0, 'The target is struck with slender, whiplike vines to inflict damage.'),
('Razor Leaf', 'grass', 'physical', 55, 95, 25, 0, 'Sharp-edged leaves are launched to slash at the opposing Pokémon. Critical hits land more easily.'),
('Solar Beam', 'grass', 'special', 120, 100, 10, 0, 'A two-turn attack. The user gathers light, then blasts a bundled beam on the next turn.'),
('Giga Drain', 'grass', 'special', 75, 100, 10, 0, 'A nutrient-draining attack. The user''s HP is restored by half the damage taken by the target.');

UPDATE moves SET crit_ratio = 1 WHERE name = 'Razor Leaf';
UPDATE moves SET drain_percent = 50 WHERE name = 'Giga Drain';

-- =====================================================
-- ICE TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Ice Beam', 'ice', 'special', 90, 100, 10, 0, 'The target is struck with an icy-cold beam of energy. This may also leave the target frozen.'),
('Blizzard', 'ice', 'special', 110, 70, 5, 0, 'A howling blizzard is summoned to strike opposing Pokémon. This may also leave the opposing Pokémon frozen.'),
('Ice Shard', 'ice', 'physical', 40, 100, 30, 1, 'The user flash-freezes chunks of ice and hurls them at the target. This move always goes first.');

UPDATE moves SET secondary_effect = '{"chance": 10, "status_inflict": {"status": "freeze", "chance": 10}}'::jsonb WHERE name = 'Ice Beam';
UPDATE moves SET secondary_effect = '{"chance": 10, "status_inflict": {"status": "freeze", "chance": 10}}'::jsonb WHERE name = 'Blizzard';

-- =====================================================
-- FIGHTING TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Karate Chop', 'fighting', 'physical', 50, 100, 25, 0, 'The target is attacked with a sharp chop. Critical hits land more easily.'),
('Low Kick', 'fighting', 'physical', 50, 100, 20, 0, 'A powerful low kick that makes the target fall over. The heavier the target, the greater the move''s power.'),
('Close Combat', 'fighting', 'physical', 120, 100, 5, 0, 'The user fights the target up close without guarding itself. This also lowers the user''s Defense and Sp. Def stats.'),
('Mach Punch', 'fighting', 'physical', 40, 100, 30, 1, 'The user throws a punch at blinding speed. This move always goes first.');

UPDATE moves SET crit_ratio = 1 WHERE name = 'Karate Chop';
UPDATE moves SET stat_changes = '[{"stat": "defense", "stages": -1, "target": "self"}, {"stat": "special_defense", "stages": -1, "target": "self"}]'::jsonb WHERE name = 'Close Combat';

-- =====================================================
-- POISON TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Poison Sting', 'poison', 'physical', 15, 100, 35, 0, 'The user stabs the target with a poisonous stinger. This may also poison the target.'),
('Sludge Bomb', 'poison', 'special', 90, 100, 10, 0, 'Unsanitary sludge is hurled at the target. This may also poison the target.'),
('Toxic', 'poison', 'status', NULL, 90, 10, 0, 'A move that leaves the target badly poisoned. Its poison damage worsens every turn.');

UPDATE moves SET secondary_effect = '{"chance": 30, "status_inflict": {"status": "poison", "chance": 30}}'::jsonb WHERE name = 'Poison Sting';
UPDATE moves SET secondary_effect = '{"chance": 30, "status_inflict": {"status": "poison", "chance": 30}}'::jsonb WHERE name = 'Sludge Bomb';
UPDATE moves SET status_inflict = '{"status": "badly_poison", "chance": 100}'::jsonb WHERE name = 'Toxic';

-- =====================================================
-- GROUND TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Dig', 'ground', 'physical', 80, 100, 10, 0, 'The user burrows, then attacks on the next turn.'),
('Earthquake', 'ground', 'physical', 100, 100, 10, 0, 'The user sets off an earthquake that strikes every Pokémon around it.'),
('Earth Power', 'ground', 'special', 90, 100, 10, 0, 'The user makes the ground under the target erupt with power. This may also lower the target''s Sp. Def stat.');

UPDATE moves SET secondary_effect = '{"chance": 10, "stat_changes": [{"stat": "special_defense", "stages": -1, "target": "opponent"}]}'::jsonb WHERE name = 'Earth Power';

-- =====================================================
-- FLYING TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Gust', 'flying', 'special', 40, 100, 35, 0, 'A gust of wind is whipped up by wings and launched at the target to inflict damage.'),
('Wing Attack', 'flying', 'physical', 60, 100, 35, 0, 'The target is struck with large, imposing wings spread wide to inflict damage.'),
('Aerial Ace', 'flying', 'physical', 60, 0, 20, 0, 'The user confounds the target with speed, then slashes. This attack never misses.'),
('Brave Bird', 'flying', 'physical', 120, 100, 15, 0, 'The user tucks in its wings and charges from a low altitude. This also damages the user quite a lot.');

UPDATE moves SET recoil_percent = 33 WHERE name = 'Brave Bird';

-- =====================================================
-- PSYCHIC TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Confusion', 'psychic', 'special', 50, 100, 25, 0, 'The target is hit by a weak telekinetic force. This may also confuse the target.'),
('Psychic', 'psychic', 'special', 90, 100, 10, 0, 'The target is hit by a strong telekinetic force. This may also lower the target''s Sp. Def stat.'),
('Future Sight', 'psychic', 'special', 120, 100, 10, 0, 'Two turns after this move is used, a hunk of psychic energy attacks the target.');

UPDATE moves SET secondary_effect = '{"chance": 10, "stat_changes": [{"stat": "special_defense", "stages": -1, "target": "opponent"}]}'::jsonb WHERE name = 'Psychic';

-- =====================================================
-- BUG TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Bug Bite', 'bug', 'physical', 60, 100, 20, 0, 'The user bites the target. If the target is holding a Berry, the user eats it and gains its effect.'),
('X-Scissor', 'bug', 'physical', 80, 100, 15, 0, 'The user slashes at the target by crossing its scythes or claws as if they were a pair of scissors.'),
('U-turn', 'bug', 'physical', 70, 100, 20, 0, 'After making its attack, the user rushes back to switch places with a party Pokémon in waiting.');

-- =====================================================
-- ROCK TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Rock Throw', 'rock', 'physical', 50, 90, 15, 0, 'The user picks up and throws a small rock at the target to attack.'),
('Rock Slide', 'rock', 'physical', 75, 90, 10, 0, 'Large boulders are hurled at the opposing Pokémon to inflict damage. This may also make the opposing Pokémon flinch.'),
('Stone Edge', 'rock', 'physical', 100, 80, 5, 0, 'The user stabs the target from below with sharpened stones. Critical hits land more easily.'),
('Stealth Rock', 'rock', 'status', NULL, 0, 20, 0, 'The user lays a trap of levitating stones around the opposing team. The trap hurts opposing Pokémon that switch into battle.');

UPDATE moves SET secondary_effect = '{"chance": 30, "flinch_chance": 30}'::jsonb WHERE name = 'Rock Slide';
UPDATE moves SET crit_ratio = 1 WHERE name = 'Stone Edge';
UPDATE moves SET entry_hazard = '{"hazard_type": "stealth_rock", "layers": 1}'::jsonb WHERE name = 'Stealth Rock';

-- =====================================================
-- GHOST TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Lick', 'ghost', 'physical', 30, 100, 30, 0, 'The target is licked with a long tongue, causing damage. This may also leave the target with paralysis.'),
('Shadow Ball', 'ghost', 'special', 80, 100, 15, 0, 'The user hurls a shadowy blob at the target. This may also lower the target''s Sp. Def stat.'),
('Shadow Claw', 'ghost', 'physical', 70, 100, 15, 0, 'The user slashes with a sharp claw made from shadows. Critical hits land more easily.');

UPDATE moves SET secondary_effect = '{"chance": 30, "status_inflict": {"status": "paralysis", "chance": 30}}'::jsonb WHERE name = 'Lick';
UPDATE moves SET secondary_effect = '{"chance": 20, "stat_changes": [{"stat": "special_defense", "stages": -1, "target": "opponent"}]}'::jsonb WHERE name = 'Shadow Ball';
UPDATE moves SET crit_ratio = 1 WHERE name = 'Shadow Claw';

-- =====================================================
-- DRAGON TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Dragon Rage', 'dragon', 'special', 40, 100, 10, 0, 'This attack hits the target with a shock wave of pure rage. This attack always inflicts 40 HP damage.'),
('Dragon Claw', 'dragon', 'physical', 80, 100, 15, 0, 'The user slashes the target with huge, sharp claws.'),
('Outrage', 'dragon', 'physical', 120, 100, 10, 0, 'The user rampages and attacks for two to three turns. The user then becomes confused.');

-- =====================================================
-- DARK TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Bite', 'dark', 'physical', 60, 100, 25, 0, 'The target is bitten with viciously sharp fangs. This may also make the target flinch.'),
('Crunch', 'dark', 'physical', 80, 100, 15, 0, 'The user crunches up the target with sharp fangs. This may also lower the target''s Defense stat.'),
('Sucker Punch', 'dark', 'physical', 70, 100, 5, 1, 'This move enables the user to attack first. This move fails if the target is not readying an attack.');

UPDATE moves SET secondary_effect = '{"chance": 30, "flinch_chance": 30}'::jsonb WHERE name = 'Bite';
UPDATE moves SET secondary_effect = '{"chance": 20, "stat_changes": [{"stat": "defense", "stages": -1, "target": "opponent"}]}'::jsonb WHERE name = 'Crunch';

-- =====================================================
-- STEEL TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Metal Claw', 'steel', 'physical', 50, 95, 35, 0, 'The target is raked with steel claws. This may also raise the user''s Attack stat.'),
('Iron Head', 'steel', 'physical', 80, 100, 15, 0, 'The user slams the target with its steel-hard head. This may also make the target flinch.'),
('Flash Cannon', 'steel', 'special', 80, 100, 10, 0, 'The user gathers all its light energy and releases it at once. This may also lower the target''s Sp. Def stat.');

UPDATE moves SET secondary_effect = '{"chance": 10, "stat_changes": [{"stat": "attack", "stages": 1, "target": "self"}]}'::jsonb WHERE name = 'Metal Claw';
UPDATE moves SET secondary_effect = '{"chance": 30, "flinch_chance": 30}'::jsonb WHERE name = 'Iron Head';
UPDATE moves SET secondary_effect = '{"chance": 10, "stat_changes": [{"stat": "special_defense", "stages": -1, "target": "opponent"}]}'::jsonb WHERE name = 'Flash Cannon';

-- =====================================================
-- FAIRY TYPE MOVES
-- =====================================================
INSERT INTO moves (name, type, category, power, accuracy, pp, priority, description) VALUES
('Fairy Wind', 'fairy', 'special', 40, 100, 30, 0, 'The user stirs up a fairy wind and strikes the target with it.'),
('Dazzling Gleam', 'fairy', 'special', 80, 100, 10, 0, 'The user damages opposing Pokémon by emitting a powerful flash.'),
('Play Rough', 'fairy', 'physical', 90, 90, 10, 0, 'The user plays rough with the target and attacks it. This may also lower the target''s Attack stat.');

UPDATE moves SET secondary_effect = '{"chance": 10, "stat_changes": [{"stat": "attack", "stages": -1, "target": "opponent"}]}'::jsonb WHERE name = 'Play Rough';
