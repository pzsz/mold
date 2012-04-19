#ifndef BBRIDGE_H_
#define BBRIDGE_H_

#ifdef __cplusplus
extern "C" {
#endif

  typedef struct _BB_World BB_World;

  typedef struct _BB_StaticMesh BB_StaticMesh;


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


#ifdef __cplusplus
}
#endif

#endif
