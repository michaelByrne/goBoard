CREATE TABLE theme
(
  id      serial PRIMARY KEY,
  name    text NOT NULL CHECK(name <> ''),
  value   text,
  main    bool NOT NULL DEFAULT false,
  UNIQUE(name)
);
INSERT INTO theme (main,name,value) VALUES (true,'blue','a:9:{s:4:"font";s:37:"verdana, helvetica, arial, sans-serif";s:8:"fontsize";s:3:"1.1";s:4:"body";s:7:"#333333";s:4:"even";s:7:"#c3dae4";s:3:"odd";s:7:"#acccdb";s:2:"me";s:7:"#82b3c9";s:5:"hover";s:7:"#82b3c9";s:7:"readbar";s:7:"#3488ab";s:4:"menu";s:7:"#555555";}');
INSERT INTO theme (name,value) VALUES ('simple','a:9:{s:4:"font";s:37:"verdana, helvetica, arial, sans-serif";s:8:"fontsize";s:3:"1.1";s:4:"body";s:7:"#ffffff";s:4:"even";s:7:"#cccccc";s:3:"odd";s:7:"#eeeeee";s:2:"me";s:7:"#82b3c9";s:5:"hover";s:7:"#82b3c9";s:7:"readbar";s:7:"#82b3c9";s:4:"menu";s:7:"#555555";}');
INSERT INTO theme (name,value) VALUES ('gray','a:9:{s:4:"font";s:37:"verdana, helvetica, arial, sans-serif";s:8:"fontsize";s:3:"1.1";s:4:"body";s:7:"#555555";s:4:"even";s:7:"#d7d7d7";s:3:"odd";s:7:"#c9c9c9";s:2:"me";s:7:"#adadad";s:5:"hover";s:7:"#adadad";s:7:"readbar";s:7:"#333333";s:4:"menu";s:7:"#555555";}');
INSERT INTO theme (name,value) VALUES ('white','a:9:{s:4:"font";s:37:"verdana, helvetica, arial, sans-serif";s:8:"fontsize";s:3:"1.1";s:4:"body";s:7:"#ffffff";s:4:"even";s:7:"#cccccc";s:3:"odd";s:7:"#eeeeee";s:2:"me";s:7:"#999999";s:5:"hover";s:7:"#999999";s:7:"readbar";s:7:"#666666";s:4:"menu";s:7:"#555555";}');
INSERT INTO theme (name,value) VALUES ('black','a:9:{s:4:"font";s:37:"verdana, helvetica, arial, sans-serif";s:8:"fontsize";s:3:"1.1";s:4:"body";s:7:"#000000";s:4:"even";s:7:"#bbbbbb";s:3:"odd";s:7:"#dddddd";s:2:"me";s:7:"#666666";s:5:"hover";s:7:"#666666";s:7:"readbar";s:7:"#555555";s:4:"menu";s:7:"#000000";}');
INSERT INTO theme (name,value) VALUES ('purple','a:9:{s:4:"font";s:37:"verdana, helvetica, arial, sans-serif";s:8:"fontsize";s:3:"1.1";s:4:"body";s:7:"#333333";s:4:"even";s:7:"#bebde9";s:3:"odd";s:7:"#a6a5e1";s:2:"me";s:7:"#7978d2";s:5:"hover";s:7:"#7978d2";s:7:"readbar";s:7:"#5553ae";s:4:"menu";s:7:"#555555";}');
INSERT INTO theme (name,value) VALUES ('green','a:9:{s:4:"font";s:37:"verdana, helvetica, arial, sans-serif";s:8:"fontsize";s:3:"1.1";s:4:"body";s:7:"#333333";s:4:"even";s:7:"#d4f0be";s:3:"odd";s:7:"#c5eba7";s:2:"me";s:7:"#a8e07a";s:5:"hover";s:7:"#a8e07a";s:7:"readbar";s:7:"#3e8c00";s:4:"menu";s:7:"#555555";}');
INSERT INTO theme (name,value) VALUES ('orange','a:9:{s:4:"font";s:37:"verdana, helvetica, arial, sans-serif";s:8:"fontsize";s:3:"1.1";s:4:"body";s:7:"#333333";s:4:"even";s:7:"#e0c18b";s:3:"odd";s:7:"#dbb878";s:2:"me";s:7:"#d1a453";s:5:"hover";s:7:"#d1a453";s:7:"readbar";s:7:"#a36d00";s:4:"menu";s:7:"#555555";}');
INSERT INTO theme (name,value) VALUES ('red','a:9:{s:4:"font";s:37:"verdana, helvetica, arial, sans-serif";s:8:"fontsize";s:3:"1.1";s:4:"body";s:7:"#333333";s:4:"even";s:7:"#a22626";s:3:"odd";s:7:"#ae2929";s:2:"me";s:7:"#7e0101";s:5:"hover";s:7:"#7e0101";s:7:"readbar";s:7:"#111111";s:4:"menu";s:7:"#555555";}');
INSERT INTO theme (name,value) VALUES ('halloween','a:9:{s:4:"font";s:37:"verdana, helvetica, arial, sans-serif";s:8:"fontsize";s:3:"1.1";s:4:"body";s:7:"#333333";s:4:"even";s:7:"#eaa61e";s:3:"odd";s:7:"#f4b028";s:2:"me";s:7:"#000000";s:5:"hover";s:7:"#000000";s:7:"readbar";s:7:"#a36d00";s:4:"menu";s:7:"#555555";}');
