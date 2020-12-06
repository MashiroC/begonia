var a = {
    "namespace": "begonia.func.Echo",
    "type": "record",
    "name": "In",
    "fields": [
        {"name": "i1", "type": "int"}
        , {"name": "i2", "type": "int"}
        , {"name": "i3", "type": "int"}
        , {"name": "i4", "type": "int"}
        , {"name": "i5", "type": "long"}
        , {"name": "f1", "type": "float"}
        , {"name": "f2", "type": "double"}
        , {"name": "ok", "type": "boolean"}
        , {"name": "str", "type": "string"}
        , {
            "name": "s1", "type": {
                "type": "array",
                "items": "int"
            }
        }
        , {
            "name": "s2", "type": {
                "type": "array",
                "items": "string"
            }
        }
        , {"name": "s6", "type": "bytes"}
        , {
            "name": "st", "type": {
                "type": "record",
                "fType": "TestStruct",
                "fields": [{"name": "I1", "type": "int"}
                    , {"name": "I2", "type": "int"}
                    , {"name": "I3", "type": "int"}
                    , {"name": "I4", "type": "int"}
                    , {"name": "I5", "type": "long"}
                    , {"name": "Str", "type": "string"}
                    , {
                        "name": "S1", "type": {
                            "type": "array",
                            "items": "int"
                        }
                    }
                    , {
                        "name": "S2", "type": {
                            "type": "array",
                            "items": "string"
                        }
                    }
                    , {
                        "name": "hello", "type": {
                            "type": "record",
                            "fType": "TestStruct2",
                            "fields": [{"name": "b1", "type": "bytes"}
                                , {"name": "b2", "type": "bytes"}

                            ]
                        }
                    }
                    , {
                        "name": "Test3", "type": {
                            "type": "record",
                            "fType": "TestStruct2",
                            "fields": [{"name": "b1", "type": "bytes"}
                                , {"name": "b2", "type": "bytes"}

                            ]
                        }
                    }
                    , {"name": "Map1", "type": {"type": "map", "values": "string"}}
                    , {
                        "name": "Map2", "type": {
                            "type": "map", "values": {
                                "type": "array",
                                "items": "int"
                            }
                        }
                    }

                ]
            }
        }
        , {"name": "m1", "type": {"type": "map", "values": "string"}}
        , {"name": "m2", "type": {"type": "map", "values": "int"}}
        , {
            "name": "m3", "type": {
                "type": "map", "values": {
                    "type": "record",
                    "fType": "TestStruct",
                    "fields": [{"name": "I1", "type": "int"}
                        , {"name": "I2", "type": "int"}
                        , {"name": "I3", "type": "int"}
                        , {"name": "I4", "type": "int"}
                        , {"name": "I5", "type": "long"}
                        , {"name": "Str", "type": "string"}
                        , {
                            "name": "S1", "type": {
                                "type": "array",
                                "items": "int"
                            }
                        }
                        , {
                            "name": "S2", "type": {
                                "type": "array",
                                "items": "string"
                            }
                        }
                        , {
                            "name": "hello", "type": {
                                "type": "record",
                                "fType": "TestStruct2",
                                "fields": [{"name": "b1", "type": "bytes"}
                                    , {"name": "b2", "type": "bytes"}

                                ]
                            }
                        }
                        , {
                            "name": "Test3", "type": {
                                "type": "record",
                                "fType": "TestStruct2",
                                "fields": [{"name": "b1", "type": "bytes"}
                                    , {"name": "b2", "type": "bytes"}

                                ]
                            }
                        }
                        , {"name": "Map1", "type": {"type": "map", "values": "string"}}
                        , {
                            "name": "Map2", "type": {
                                "type": "map", "values": {
                                    "type": "array",
                                    "items": "int"
                                }
                            }
                        }

                    ]
                }
            }
        }

    ]
}