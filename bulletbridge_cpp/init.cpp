#include "bbridge.h"
#include "btBulletDynamicsCommon.h"

struct _BB_World {
  btDiscreteDynamicsWorld          *dynamicsWorld;
  btBroadphaseInterface	           *broadphase;
  btCollisionDispatcher	           *dispatcher;
  btConstraintSolver               *solver;
  btDefaultCollisionConfiguration  *collisionConfiguration;
};

struct _BB_StaticMesh {
  BB_World                   *world;
  btCollisionObject          *collisionObject;
  btTriangleIndexVertexArray *mesh;
  btBvhTriangleMeshShape     *shape;
};


BB_World* BB_NewWorld() {
  BB_World* ret = new BB_World();

  ret->collisionConfiguration = new btDefaultCollisionConfiguration();
  ret->dispatcher = new btCollisionDispatcher(ret->collisionConfiguration);

  btVector3 worldAabbMin(-100, -100, -100);
  btVector3 worldAabbMax(100, 100, 100);

  ret->broadphase = new btAxisSweep3(worldAabbMin, worldAabbMax);

  ret->solver = new btSequentialImpulseConstraintSolver();
  ret->dynamicsWorld = 
    new btDiscreteDynamicsWorld(
				ret->dispatcher, 
				ret->broadphase,
				ret->solver, 
				ret->collisionConfiguration);  
  return ret;
}


void BB_DestroyWorld(BB_World* world) {
  delete world->dynamicsWorld;
  delete world->solver;
  delete world->broadphase;
  delete world->dispatcher;
  delete world->collisionConfiguration;
  delete world;
}

void BB_ProcessWorld(BB_World* world, float time_step) {
  world->dynamicsWorld->stepSimulation(time_step, 10);
}


BB_StaticMesh* BB_NewStaticMesh(BB_World* world,
				int vertex_size, 
				char *vertex_buffer, 
				int vertex_buffer_size, 
				char *indice_buffer, 
				int indice_buffer_size) {
  
  BB_StaticMesh* mesh = new BB_StaticMesh();
  mesh->world = world;

  mesh->collisionObject = new btCollisionObject();

  btIndexedMesh imesh;
  imesh.m_numTriangles = indice_buffer_size/(2*3);
  imesh.m_triangleIndexBase = (const unsigned char *)indice_buffer;
  imesh.m_triangleIndexStride = 6;
  imesh.m_numVertices = vertex_buffer_size/vertex_size;
  imesh.m_vertexBase = (const unsigned char *)vertex_buffer;
  imesh.m_vertexStride = vertex_size;

  mesh->mesh = new btTriangleIndexVertexArray();
  mesh->mesh->addIndexedMesh(imesh, PHY_SHORT);

  mesh->shape = new btBvhTriangleMeshShape(mesh->mesh, true);

  mesh->collisionObject->setCollisionShape(mesh->shape);

  btTransform groundTransform;
  groundTransform.setIdentity();
  groundTransform.setOrigin(btVector3(0,0,0));
  mesh->collisionObject->setWorldTransform(groundTransform);
  
  world->dynamicsWorld->addCollisionObject(mesh->collisionObject);

  return mesh;
}

void BB_DestroyStaticMesh(BB_StaticMesh* mesh) {

}
