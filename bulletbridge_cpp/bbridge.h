#ifndef BBRIDGE_H_
#define BBRIDGE_H_

#ifdef __cplusplus
extern "C" {
#endif
  typedef struct _BB_Vector3 {
    float x;
    float y;
    float z;
  } BB_Vector3;

  typedef struct _BB_World BB_World;

  typedef struct _BB_StaticMesh BB_StaticMesh;

  typedef struct _BB_RBody BB_RBody;

  typedef struct _BB_CharacterControler BB_CharacterControler;

  typedef struct _BB_CShape {
    void *ptr;
  }BB_CShape;

  extern BB_World* BB_NewWorld();

  extern void BB_DestroyWorld(BB_World* world);

  extern void BB_ProcessWorld(BB_World* world, float time_step);

  extern BB_StaticMesh* BB_NewStaticMesh(BB_World* world,
					 int vertex_size, 
					 char *vertex_buffer, 
					 int vertex_buffer_size, 
					 char *indice_buffer, 
					 int indice_buffer_size);

  extern void BB_DestroyStaticMesh(BB_StaticMesh* mesh);

  
  extern BB_CShape BB_NewCShapeSphere(float radius);

  extern BB_RBody* BB_NewRBody(BB_World* world, 
			       BB_CShape shape, 
			       float mass, BB_Vector3 pos);

  extern void BB_DestroyRBody();

  extern void BB_SetPositionRBody(BB_RBody* rbody, BB_Vector3 pos);

  extern BB_Vector3 BB_GetPositionRBody(BB_RBody* rbody);

  extern BB_CharacterControler* BB_NewCharacterControler(BB_World* world, float height, float width, BB_Vector3 pos);

  extern void BB_DestroyCharacterControler(BB_CharacterControler* character);

  extern BB_Vector3 BB_GetPositionCharacterControler(BB_CharacterControler* characted);
  extern void BB_SetWalkDirection(BB_CharacterControler* characted, BB_Vector3 walk);



#ifdef __cplusplus
}
#endif

#endif
