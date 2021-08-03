export interface PlayerClass {
  id: number;
  name: string;
  power_type: PowerType;
  specializations: Specialization[];
}

export interface PowerType {
  name: string;
  id: number;
}

export interface Specialization {
  name: string;
  id: number;
}
